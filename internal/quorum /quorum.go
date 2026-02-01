package quorum

func Calculate(total int) int {
	return total/2 + 1
}

func IsReached(responses, total int) bool {
	return responses >= Calculate(total)
}
