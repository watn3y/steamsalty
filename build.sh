#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME=$(basename "$PWD")

if [[ ! -f go.mod ]]; then
    echo "go.mod not found in $PWD" >&2
    exit 1
fi

# Parse the "go X.Y[.Z]" directive from go.mod
GO_VERSION=$(awk '/^go [0-9]/ {print $2; exit}' go.mod)
if [[ -z "$GO_VERSION" ]]; then
    echo "Could not parse Go version from go.mod" >&2
    exit 1
fi

echo "Project:    $PROJECT_NAME"
echo "Go version: $GO_VERSION"
echo "Image:      golang:$GO_VERSION"

architectures=(
    "linux/amd64"
    "linux/386"
    "linux/arm64"
    "linux/arm/v7"
    "linux/riscv64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

mkdir -p build

read -r -d '' BUILD_SCRIPT <<'INNER' || true
PROJECT_NAME="$1"
shift

export GOCACHE=/tmp/go-build
export GOMODCACHE=/tmp/go-mod

failed=0
for arch in "$@"; do
    os="${arch%%/*}"
    rest="${arch#*/}"
    if [[ "$rest" == *"/"* ]]; then
        arch_type="${rest%%/*}"
        arm_version="${rest#*/}"
    else
        arch_type="$rest"
        arm_version=""
    fi

    suffix="$os-$arch_type"
    [[ -n "$arm_version" ]] && suffix="$suffix-$arm_version"

    ext=""
    [[ "$os" == "windows" ]] && ext=".exe"

    output_file="build/$PROJECT_NAME-$suffix$ext"

    env_args=("GOOS=$os" "GOARCH=$arch_type")
    if [[ "$arch_type" == "arm" && -n "$arm_version" ]]; then
        env_args+=("GOARM=${arm_version#v}")
    fi

    echo ">>> Building $arch"
    if env "${env_args[@]}" go build -o "$output_file" .; then
        if tar -czf "$output_file.tar.gz" -C build "$(basename "$output_file")"; then
            rm "$output_file"
            echo "    OK:   $output_file.tar.gz"
        else
            echo "    FAIL: archive for $output_file" >&2
            failed=1
        fi
    else
        echo "    FAIL: build $arch" >&2
        failed=1
    fi
done

exit $failed
INNER

docker run --rm \
-v "$PWD:/app" \
-w /app \
--user "$(id -u):$(id -g)" \
-e HOME=/tmp \
"golang:$GO_VERSION" \
bash -c "$BUILD_SCRIPT" _ "$PROJECT_NAME" "${architectures[@]}"
