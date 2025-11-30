package ssm

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("AWS_REGION", "eu-west-2")
}

func TestSsm(t *testing.T) {
	t.Run("can load all params", func(t *testing.T) {
		t.Skip("skip test relying on AWS --profile joe calling live API")
		sc := NewClient()

		sc.LoadParameterValues()

		list := []Parameter{
			sc.GoogleClientEmail,
			sc.GooglePrivateKey,
		}

		missing := []string{}
		for _, p := range list {
			if p.Value == "" {
				missing = append(missing, p.Name)
			}
		}

		if len(missing) > 0 {
			t.Errorf("Failed to load params:%v\n", missing)
		}
	})
}
