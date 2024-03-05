package notion

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStudyStepIsExpired(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		step           StudyStep
		expectedResult bool
		desc           string
	}

	for _, scneario := range []testCase{
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2024, 02, 01, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: false,
			desc:           "not expired, deadline in a month",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: true,
			desc:           "expired, exactly the same day",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2024, 01, 02, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: false,
			desc:           "not expired, deadline in a day",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: true,
			desc:           "expired, deadline expired a month ago",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Time{},
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: false,
			desc:           "not expired, deadline is zero",
		},
	} {
		t.Run(scneario.desc, func(t *testing.T) {
			now = time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)
			assert.Equal(scneario.step.IsExpired(), scneario.expectedResult)
		})
	}

}

func TestStudyStepExpireSoon(t *testing.T) {
	assert := assert.New(t)

	type testCase struct {
		step           StudyStep
		expectedResult bool
		desc           string
	}

	for _, scneario := range []testCase{
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Finished,
				Notes:     "",
				Deadline:  time.Date(2024, 02, 01, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: false,
			desc:           "not Expire Soon, status is Finished",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: false,
			desc:           "not Expire Soon, is already expired",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2024, 01, 15, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: true,
			desc:           "Expire Soon, expire exactly a day before deadline",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2024, 01, 17, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: false,
			desc:           "not Expire Soon, deadline expires in a day",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Date(2024, 01, 16, 0, 0, 0, 0, time.UTC),
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: true,
			desc:           "Expire Soon, expire exactly at deadline day",
		},
		{
			step: StudyStep{
				ID:        "test123",
				Name:      "name1",
				Category:  "Category1",
				Link:      "www.test.com",
				Status:    Blocked,
				Notes:     "",
				Deadline:  time.Time{},
				Priority:  Medium,
				CreatedAt: time.Date(2023, 12, 01, 0, 0, 0, 0, time.UTC),
				Type:      []string{"type1", "type2"},
				StartedAt: time.Date(2023, 20, 01, 0, 0, 0, 0, time.UTC),
				Points:    5,
			},
			expectedResult: false,
			desc:           "not Expire Soon, deadline is zero",
		},
	} {
		t.Run(scneario.desc, func(t *testing.T) {
			now = time.Date(2024, 01, 01, 0, 0, 0, 0, time.UTC)
			expiredDays = 15
			assert.Equal(scneario.step.ExpireSoon(), scneario.expectedResult)
		})
	}
}
