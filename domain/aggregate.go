package domain

import "github.com/blucv2022/crowdstats/models"

func (smr *Summarizer) aggregateK(digests []*models.DataDigest) *models.DataDigest {
	// n := len(digests)

	// if n == 0 {
	// 	return nil
	// }

	//TODO: modify below to use "divide and conquer" as specified in README
	//      and add tests
	// if n == 1 {
	// 	return digests[0]
	// }
	// resDigest := digests[0]
	// for i := 1; i < n; i++ {
	// 	resDigest = smr.aggregateTwo(resDigest, digests[i])
	// }
	// return resDigest

	k := len(digests)
	jump := 1

	for jump < k {
		for i := 0; i < k-jump; i += jump * 2 {
			digests[i] = smr.aggregateTwo(digests[i], digests[i+jump])
		}
		jump *= 2
	}

	if k > 0 {
		return digests[0]
	} else {
		return nil
	}
}

func (smr *Summarizer) aggregateTwo(a, b *models.DataDigest) *models.DataDigest {
	ageCounts := [201]int64{}
	agePersonName := [201]string{}

	for i := 0; i <= 200; i++ {
		ageCounts[i] = a.AgeCounts[i] + b.AgeCounts[i]

		if len(a.AgePersonName[i]) > 0 {
			agePersonName[i] = a.AgePersonName[i]
		} else if len(b.AgePersonName[i]) > 0 {
			agePersonName[i] = b.AgePersonName[i]
		}
	}

	return &models.DataDigest{
		TotalAgeCounts: a.TotalAgeCounts + b.TotalAgeCounts,
		AgeCounts:      ageCounts,
		AgePersonName:  agePersonName,
	}
}
