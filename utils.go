package main

func get_last_n_characters_in_string(in_string string, n int) string {
	runes := []rune(in_string)
	start := int_max(len(runes)-n, 0)
	return string(runes[start:])
}

func get_first_n_characters_in_string(in_string string, n int) string {
	runes := []rune(in_string)
	end := int_min(n, len(runes))
	return string(runes[:end])
}

func remove_last_letter(text string) string {
	new_length := int_max(0, len(text)-1)
	return get_first_n_characters_in_string(text, new_length)
}

func remove_first_letter(text string) string {
	new_length := int_max(0, len(text)-1)
	return get_last_n_characters_in_string(text, new_length)
}

func int_max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func int_min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
