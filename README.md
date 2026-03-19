# 🛡️ Sentinel

**High-performance, single-binary Vercel/Netlify alternative for self-hosters.**

[![Go Report Card](https://goreportcard.com/badge/github.com/Rav-xyl/sentinel)](https://goreportcard.com/report/github.com/Rav-xyl/sentinel)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Sentinel is a lightweight, zero-dependency deployment engine written in Go. It turns any bare-metal server or $5/month VPS into a fully automated hosting platform with zero-config HTTPS, dynamic routing, and built-in CI/CD.

## 🚀 The Vision: A "SaaS Killer"
Developers are experiencing severe "SaaS Fatigue," paying $20/mo per seat for basic hosting. Sentinel gives you the same "Push to Deploy" Developer Experience (DX) on your own hardware, without requiring a PhD in Kubernetes or complex NGINX configurations.

### Key Differentiators
- **Single Binary:** Just download and run. No Docker footprint required for the proxy core.
- **SQLite State:** Zero external databases to manage. Routing rules and deployment history are stored locally.
- **Zero-Config SSL:** Automatic provisioning and renewal of Let's Encrypt certificates based dynamically on your routing table.
- **Sub-10ms Overhead:** Built with aggressive connection pooling and strict timeouts to ensure edge-level performance.

---

## 🛠️ Architecture

Sentinel operates as a high-performance gateway:

1. **Proxy Core:** A concurrent HTTP/HTTPS reverse proxy utilizing `net/http/httputil`.
2. **Dynamic Router:** Backed by SQLite, allowing routing rules to be updated without server restarts ("Zero-Restart" design).
3. **Security Layer:** Built-in per-IP token bucket rate limiting (DDoS protection) and Slowloris mitigation.
4. **Control Plane (In Development):** A React/Vite dashboard for managing deployments via GitHub webhooks.

---

## 📦 Quick Start

### 1. Installation
Ensure you have Go 1.22+ installed.
```bash
git clone https://github.com/Rav-xyl/sentinel.git
cd sentinel
go build -o sentinel
```

### 2. Run the Engine
Requires root privileges to bind to ports 80 and 443.
```bash
sudo ./sentinel
```

*Note: On first run, Sentinel will generate `sentinel.db` and start listening for traffic.*

---

## 🤝 Contributing
We welcome contributions! Please see our [CONTRIBUTING.md](CONTRIBUTING.md) for details on our zero-emoji policy, testing mandates, and architectural guidelines.

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
