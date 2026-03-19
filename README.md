# 🛡️ Sentinel

**High-performance, single-binary Vercel/Netlify alternative for self-hosters.**

Sentinel is a lightweight, zero-dependency deployment engine written in Go. It turns any bare-metal server or $5/month VPS into a fully automated hosting platform with zero-config HTTPS, dynamic routing, and built-in CI/CD.

## 🚀 Why Sentinel?
- **SaaS Killer:** Stop paying $20/mo per seat for basic hosting. Sentinel gives you the same "Push to Deploy" experience on your own hardware.
- **Single Binary:** No complex NGINX configs or Docker dependencies. Just one binary to rule them all.
- **Zero-Config SSL:** Automatic provisioning and renewal of Let's Encrypt certificates.
- **Edge Performance:** Built on Go's lightning-fast `net/http` stack for high-throughput request proxying.

## 🛠️ Core Architecture
Sentinel consists of three primary components:
1. **The Proxy Core:** A concurrent HTTP/HTTPS reverse proxy that routes traffic based on hostnames.
2. **The Dynamic Router:** A thread-safe routing table that maps domains to local processes or containers.
3. **The Control Plane (Upcoming):** A React-powered dashboard for managing deployments, viewing logs, and configuring webhooks.

## 📦 Installation
```bash
# Clone and build
git clone https://github.com/Rav-xyl/sentinel.git
cd sentinel
go build -o sentinel
# Start the engine
./sentinel
```

## 📅 Progress Tracking
See [JOURNAL.md](./JOURNAL.md) for the daily build log and 30-day roadmap.
