package stream

func selectMime(
	container StreamContainer,
	supported map[string][]string,
	isScreen bool,
) string {
	// If TV returned nothing → fallback
	if len(supported) == 0 {
		if isScreen {
			return container.MimeCandidates()[0] // Use container's primary MIME
		}
		return "video/mpeg"
	}

	candidates := container.MimeCandidates()
	// For screen, prioritize container's candidates without forcing video/mpeg
	if isScreen {
		// No special append; use container's order
	}

	for _, cand := range candidates {
		if _, ok := supported[cand]; ok {
			return cand
		}
	}

	// Nothing matched → safe fallback
	if isScreen {
		return container.MimeCandidates()[0]
	}
	return "video/mpeg"
}
