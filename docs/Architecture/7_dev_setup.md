# 🤖 Developer Setup

## **Prerequisites**

Before you begin, ensure you have the following tools installed:

-   [Go](https://go.dev/doc/install)
-   [Node.js](https://nodejs.org/en/download/)
-   [Bun](https://bun.sh/docs/installation)
-   [Mise](https://mise.jdx.dev/getting-started.html)
-   [Temporal](https://docs.temporal.io/cli)

## **Installation Steps**

### 1. Install `slangroom-exec`

Download the appropriate `slangroom-exec` binary for your OS from the [releases page](https://github.com/dyne/slangroom-exec/releases).

Add `slangroom-exec` to PATH and make it executable:

```bash
wget https://github.com/dyne/slangroom-exec/releases/latest/download/slangroom-exec-Linux-x86_64 -O slangroom-exec
chmod +x slangroom-exec
sudo cp slangroom-exec /usr/local/bin/
```

### 2. Install `zencode-exec`

Download the appropriate `zencode-exec` binary for your OS from the [releases page](https://github.com/dyne/zenroom/releases).

> [!WARNING]
> On Mac OS, the executable is zencode-exec.command and you have to symlink it to zencode-exec

Add `zencode-exec` to PATH and make it executable:

```bash
wget https://github.com/dyne/zenroom/releases/latest/download/zencode-exec
chmod +x zencode-exec
sudo cp zencode-exec /usr/local/bin/
```

### 3. **Setup Project**

Clone the repository:

```bash
git clone https://github.com/ForkbombEu/DIDimo
```

Initialize the environment:

```bash
cd DIDimo
make didimo
```

### 4. **Start Development Server**

```bash
make dev
```
