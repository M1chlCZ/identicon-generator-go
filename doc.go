// Package identicon generates deterministic, GitHub-style avatar images from
// a seed.
//
// Render accepts a string seed and validates its configuration:
//
//	img, err := identicon.Render("alice@example.com", identicon.DefaultConfig())
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := identicon.SaveImage(img, "alice.png"); err != nil {
//		log.Fatal(err)
//	}
//
// GenerateIdenticon and GenerateIdenticonWithConfig keep the original
// integer-seed and MD5 output for existing callers.
package identicon
