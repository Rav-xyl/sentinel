# Contributing to Sentinel

First off, thank you for considering contributing to Sentinel! Our goal is to build the highest-performance, most developer-friendly self-hosting engine in the world. 

To maintain the quality and security of this project, we enforce strict contribution guidelines. By participating, you agree to abide by these rules.

## 🛡️ The Zero Trust Security Mandate
Security is our highest priority.
1. **Never commit secrets:** Never commit `.env` files, API keys, or private certificates.
2. **Review `.gitignore`:** Ensure your local development files are covered by the `.gitignore` before staging any changes.
3. **No third-party binaries:** Do not include compiled binaries or opaque blobs in your Pull Requests.

## 🧪 Testing is the Universal Language
A Pull Request without tests is an incomplete Pull Request.
1. **Prove your fix:** If you fix a bug, you must add a unit test that fails without your fix and passes with it.
2. **Benchmark performance:** If you modify the core proxy (`proxy/router.go`), you must ensure your changes do not introduce latency. Sub-10ms overhead is the standard.

## 📝 Code Style & Commits
We maintain a highly professional repository environment.
1. **Zero Emojis:** Do not use emojis in commit messages, PR titles, or source code comments.
2. **Conventional Commits:** Use standard prefixes:
   - `feat(scope): ...` for new features.
   - `fix(scope): ...` for bug fixes.
   - `perf(scope): ...` for performance improvements.
   - `docs(scope): ...` for documentation changes.
3. **Go Formatting:** All code must pass `go fmt` and `go vet` before submission.

## 🚀 How to Submit a Pull Request
1. **Fork the repository** and create a topic branch (`fix/issue-description` or `feat/feature-name`).
2. **Discuss first:** For major architectural changes, please open an Issue to discuss your design before writing code.
3. **Draft a detailed PR:** Explain the *Why*, the *How*, and the *Verification* steps you took.

Thank you for helping us build the ultimate SaaS killer.