package postgres

import (
	"PRService/internal/model"
	"PRService/internal/service"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testRepo *Repository
var testDBPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		fmt.Printf("failed to start postgres container: %s", err)
		os.Exit(1)
	}

	// 2. Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Printf("failed to get connection string: %s", err)
		os.Exit(1)
	}

	// 3. Connect to the database
	testDBPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Printf("failed to connect to database: %s", err)
		os.Exit(1)
	}
	testRepo = &Repository{c: testDBPool}

	// 4. Apply schema
	schemaPath := filepath.Join("../../mocks", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Printf("failed to read schema file: %s", err)
		os.Exit(1)
	}
	_, err = testDBPool.Exec(ctx, string(schema))
	if err != nil {
		fmt.Printf("failed to apply schema: %s", err)
		os.Exit(1)
	}

	// 5. Run tests
	code := m.Run()

	// 6. Terminate container
	if err := pgContainer.Terminate(ctx); err != nil {
		fmt.Printf("failed to terminate container: %s", err)
	}

	os.Exit(code)
}

func cleanup(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	_, err := testDBPool.Exec(ctx, "TRUNCATE teams, users, pull_requests, pull_request_reviewers RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func TestTeamExists(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	team := &model.Team{Name: "team-a", Members: []model.User{{ID: "user-1", Name: "User One"}}}
	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)

	exists, err := testRepo.TeamExists(ctx, "team-a")
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = testRepo.TeamExists(ctx, "team-b")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestCreateAndGetTeam(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	team := &model.Team{
		Name: "team-alpha",
		Members: []model.User{
			{ID: "user-1", Name: "Alice", IsActive: true, TeamName: "team-alpha"},
			{ID: "user-2", Name: "Bob", IsActive: false, TeamName: "team-alpha"},
		},
	}

	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)

	// Get the team back
	fetchedTeam, err := testRepo.GetTeam(ctx, "team-alpha")
	require.NoError(t, err)
	assert.Equal(t, team.Name, fetchedTeam.Name)
	assert.Len(t, fetchedTeam.Members, 2)
	assert.ElementsMatch(t, team.Members, fetchedTeam.Members)
}

func TestGetUserAndSetUserIsActive(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	team := &model.Team{Name: "team-a", Members: []model.User{{ID: "user-1", Name: "User One", IsActive: true}}}
	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)

	// Check existence via GetUser
	user, err := testRepo.GetUser(ctx, "user-1")
	assert.NoError(t, err)
	assert.Equal(t, model.UserID("user-1"), user.ID)

	// Check non-existence
	_, err = testRepo.GetUser(ctx, "user-nonexistent")
	assert.ErrorIs(t, err, service.ErrResourceNotFound)

	// Set inactive
	userToUpdate := &model.User{ID: "user-1", IsActive: false}
	err = testRepo.SetUserIsActive(ctx, userToUpdate)
	require.NoError(t, err)
	assert.False(t, userToUpdate.IsActive)
	assert.Equal(t, model.TeamName("team-a"), userToUpdate.TeamName) // Check that other fields are returned

	// Verify it's inactive in the DB
	fetchedTeam, err := testRepo.GetTeam(ctx, "team-a")
	require.NoError(t, err)
	assert.False(t, fetchedTeam.Members[0].IsActive)
}

func TestCreatePullRequest(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	// Setup: 1 team, 3 users
	team := &model.Team{
		Name: "team-awesome",
		Members: []model.User{
			{ID: "author-1", Name: "Author", IsActive: true, TeamName: "team-awesome"},
			{ID: "reviewer-1", Name: "Reviewer One", IsActive: true, TeamName: "team-awesome"},
			{ID: "reviewer-2", Name: "Reviewer Two", IsActive: true, TeamName: "team-awesome"},
		},
	}
	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)

	pr := &model.PullRequest{
		ID:       "pr-123",
		Name:     "My new feature",
		AuthorID: "author-1",
	}
	pr.Init()

	err = testRepo.CreatePullRequest(ctx, pr)
	require.NoError(t, err)

	// Assertions
	assert.Len(t, pr.Reviewers, 2, "Should assign 2 reviewers")
	assert.NotContains(t, pr.Reviewers, model.UserID("author-1"), "Author should not be a reviewer")

	// Get it back to verify
	fetchedPR, err := testRepo.GetPullRequest(ctx, "pr-123")
	require.NoError(t, err)
	assert.Equal(t, pr.ID, fetchedPR.ID)
	assert.Equal(t, pr.Name, fetchedPR.Name)
	assert.Equal(t, pr.AuthorID, fetchedPR.AuthorID)
	assert.Equal(t, model.PullRequestStatusOpen, fetchedPR.Status)
	assert.ElementsMatch(t, pr.Reviewers, fetchedPR.Reviewers)
}

func TestCreatePullRequest_NoAvailableReviewers(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	// Setup: 1 team, 1 user
	team := &model.Team{Name: "team-solo", Members: []model.User{{ID: "author-1", Name: "Author"}}}
	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)

	pr := &model.PullRequest{ID: "pr-456", Name: "My lonely feature", AuthorID: "author-1"}
	pr.Init()

	err = testRepo.CreatePullRequest(ctx, pr)
	assert.ErrorIs(t, err, service.ErrReviewerNotAssigned)
}

func TestGetPullRequest_NotFound(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	_, err := testRepo.GetPullRequest(ctx, "pr-nonexistent")
	assert.ErrorIs(t, err, service.ErrResourceNotFound)
}

func TestMergePullRequest(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	// Setup
	team := &model.Team{Name: "team-a", Members: []model.User{
		{ID: "author-1", Name: "Author", IsActive: true, TeamName: "team-a"},
		{ID: "reviewer-1", Name: "Reviewer One", IsActive: true, TeamName: "team-a"},
		{ID: "reviewer-2", Name: "Reviewer Two", IsActive: true, TeamName: "team-a"},
	}}
	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)
	pr := &model.PullRequest{ID: "pr-to-merge", Name: "Ready to merge", AuthorID: "author-1"}
	pr.Init()
	err = testRepo.CreatePullRequest(ctx, pr)
	require.NoError(t, err)

	// Merge
	mergedPR, err := testRepo.MergePullRequest(ctx, "pr-to-merge")
	require.NoError(t, err)

	// Assertions
	assert.Equal(t, model.PullRequestStatusMerged, mergedPR.Status)
	assert.NotZero(t, mergedPR.MergedAt)
	assert.Equal(t, pr.ID, mergedPR.ID)

	// Verify in DB
	fetchedPR, err := testRepo.GetPullRequest(ctx, "pr-to-merge")
	require.NoError(t, err)
	assert.Equal(t, model.PullRequestStatusMerged, fetchedPR.Status)
	assert.False(t, fetchedPR.MergedAt.IsZero())
}

func TestGetReviewersPRs(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	// Setup
	team := &model.Team{
		Name: "team-review",
		Members: []model.User{
			{ID: "author-1", Name: "Author One", IsActive: true, TeamName: "team-review"},
			{ID: "author-2", Name: "Author Two", IsActive: true, TeamName: "team-review"},
			{ID: "reviewer-1", Name: "The Reviewer", IsActive: true, TeamName: "team-review"},
		},
	}
	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)

	// Create two PRs where "reviewer-1" will be assigned
	pr1 := &model.PullRequest{ID: "pr-1", Name: "PR One", AuthorID: "author-1"}
	pr1.Init()
	err = testRepo.CreatePullRequest(ctx, pr1)
	require.NoError(t, err)

	pr2 := &model.PullRequest{ID: "pr-2", Name: "PR Two", AuthorID: "author-2"}
	pr2.Init()
	err = testRepo.CreatePullRequest(ctx, pr2)
	require.NoError(t, err)

	// Manually ensure reviewer-1 is on both PRs for deterministic test
	_, err = testDBPool.Exec(ctx, "DELETE FROM pull_request_reviewers")
	require.NoError(t, err)
	_, err = testDBPool.Exec(ctx, "INSERT INTO pull_request_reviewers (pull_request_id, user_id) VALUES ('pr-1', 'reviewer-1'), ('pr-2', 'reviewer-1')")
	require.NoError(t, err)

	// Test
	prs, err := testRepo.GetReviewersPRs(ctx, "reviewer-1")
	require.NoError(t, err)

	assert.Len(t, prs, 2)
	prIDs := []model.PullRequestID{prs[0].ID, prs[1].ID}
	assert.ElementsMatch(t, []model.PullRequestID{"pr-1", "pr-2"}, prIDs)
}

func TestReassignPullRequestReviewer(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	// Setup: 1 team, 4 users
	team := &model.Team{
		Name: "team-reassign",
		Members: []model.User{
			{ID: "author-1", Name: "Author"},
			{ID: "old-reviewer", Name: "Old Reviewer"},
			{ID: "other-reviewer", Name: "Other Reviewer"},
			{ID: "new-reviewer", Name: "New Reviewer"},
		},
	}
	err := testRepo.CreateTeam(ctx, team)
	require.NoError(t, err)

	// Manually create PR and reviewers for a deterministic test
	_, err = testDBPool.Exec(ctx, `INSERT INTO pull_requests (id, name, author_id, status, created_at) 
										VALUES ('pr-reassign', 'Reassign Test', 'author-1', 'OPEN', NOW())`)
	require.NoError(t, err)

	_, err = testDBPool.Exec(ctx, "INSERT INTO pull_request_reviewers (pull_request_id, user_id) VALUES ('pr-reassign', 'old-reviewer'), ('pr-reassign', 'other-reviewer')")
	require.NoError(t, err)

	// Fetch the PR to pass to the method
	prToUpdate, err := testRepo.GetPullRequest(ctx, "pr-reassign")
	require.NoError(t, err)
	require.Len(t, prToUpdate.Reviewers, 2)

	// Reassign "old-reviewer"
	err = testRepo.ReassignPullRequestReviewer(ctx, &prToUpdate, "old-reviewer", "new-reviewer")
	require.NoError(t, err)

	// Assertions
	// The storage method doesn't modify the passed-in struct's reviewers,
	// so we need to fetch it again to verify the change in the database.

	// Verify in DB
	finalPR, err := testRepo.GetPullRequest(ctx, "pr-reassign")
	require.NoError(t, err)
	assert.Len(t, finalPR.Reviewers, 2)
	assert.NotContains(t, finalPR.Reviewers, model.UserID("old-reviewer"))
	assert.Contains(t, finalPR.Reviewers, model.UserID("other-reviewer"))
	assert.Contains(t, finalPR.Reviewers, model.UserID("new-reviewer"))
}

func TestGetStatistic(t *testing.T) {
	cleanup(t)
	ctx := context.Background()

	// Setup data
	_, err := testDBPool.Exec(ctx, `
		INSERT INTO teams (name) VALUES ('team-1'), ('team-2');
		INSERT INTO users (id, name, team_name, is_active) VALUES 
			('u1', 'a', 'team-1', true), 
			('u2', 'b', 'team-1', true), 
			('u3', 'c', 'team-1', false);
		INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at) VALUES
			('pr1', 'p1', 'u1', 'OPEN', NOW() - interval '2 day', null),
			('pr2', 'p2', 'u2', 'MERGED', NOW() - interval '1 day', NOW());
	`)
	require.NoError(t, err)

	stats, err := testRepo.GetStatistic(ctx)
	require.NoError(t, err)

	assert.Equal(t, 2, stats.TotalTeams)
	assert.Equal(t, 3, stats.TotalUsers)
	assert.Equal(t, 2, stats.ActiveUsers)
	assert.Equal(t, 2, stats.TotalPRs)
	assert.Equal(t, 1, stats.OpenPRs)
	assert.Equal(t, 1, stats.MergedPRs)
	assert.InDelta(t, 86400, stats.AvgMergeTimeSeconds, 1.0, "average merge time should be around 1 day (86400s)")
}
