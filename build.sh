#!/usr/bin/env bash


# Get the current directory name as the project name
PROJECT_NAME=$(basename "$PWD")

# Change to the Go project directory (current directory)
cd "$PWD" || exit

# List of target architectures
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

# Create a build directory if it doesn't exist
mkdir -p build

# Loop over architectures and build for each one
for arch in "${architectures[@]}"; do
  os=$(echo $arch | cut -d '/' -f 1)
  arch_type=$(echo $arch | cut -d '/' -f 2)
  
  # Set the output file name using the project name, os, and arch
  output_file="build/$PROJECT_NAME-$os-$arch_type"

  # Build the app for the specific architecture
  GOARCH=$arch_type GOOS=$os go build -o "$output_file"

  # Provide feedback to the user
  if [ $? -eq 0 ]; then
    echo "Successfully built for $arch: $output_file"
    
    # Create a tar.gz archive for the individual build
    tar -czf "$output_file.tar.gz" -C build $(basename "$output_file")
    if [ $? -eq 0 ]; then
      echo "Successfully created archive: $output_file.tar.gz"
      
      # Delete the executable after archiving
      rm "$output_file"
      echo "Deleted executable: $output_file"
    else
      echo "Failed to create archive for $output_file"
    fi
  else
    echo "Failed to build for $arch"
  fi
done