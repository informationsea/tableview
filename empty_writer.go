package main

type EmptyWriter struct{}

func (w EmptyWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
