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
chmod +x slangroom-exec
echo "export PATH=\$PATH:$(pwd)/slangroom-exec" >> ~/.bashrc
```

### 2. **Setup Project**

Clone the repository:

```bash
git clone https://github.com/ForkbombEu/DIDimo
```

Initialize the environment:

```bash
cd DIDimo
make didimo
```

### 3. **Start Development Server**

```bash
make dev
```
