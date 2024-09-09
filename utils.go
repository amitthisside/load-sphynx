package main

// Helper function to calculate the greatest common divisor (GCD) of the weights
func gcd(servers []*Server) int {
	gcd := servers[0].Weight
	for _, server := range servers {
		gcd = gcdTwoNumbers(gcd, server.Weight)
	}
	return gcd
}

// Helper function to calculate the GCD of two numbers
func gcdTwoNumbers(a, b int) int {
	if b == 0 {
		return a
	}
	return gcdTwoNumbers(b, a%b)
}

// Helper function to find the maximum weight among the servers
func maxWeight(servers []*Server) int {
	max := 0
	for _, server := range servers {
		if server.Weight > max {
			max = server.Weight
		}
	}
	return max
}
