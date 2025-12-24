#!/bin/bash
set -e
# Get current version from latest git tag
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
CURRENT_VERSION=${CURRENT_VERSION#v}  # Remove 'v' prefix
# Parse version components
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"
# Determine bump type
BUMP_TYPE=${1:-patch}
case $BUMP_TYPE in
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH + 1))
        ;;
    *)
        echo "Usage: $0 [major|minor|patch]"
        exit 1
        ;;
esac
NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"
echo "Current version: v${CURRENT_VERSION}"
echo "New version: ${NEW_VERSION}"
# Check for uncommitted changes
if [[ -n $(git status --porcelain) ]]; then
    echo "Warning: You have uncommitted changes"
    git status --short
    echo ""
fi
# Create and push tag
read -p "Create and push tag ${NEW_VERSION}? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git tag -a "$NEW_VERSION" -m "Release ${NEW_VERSION}"
    git push origin "$NEW_VERSION"
    echo "✓ Tagged and pushed ${NEW_VERSION}"
else
    echo "Aborted."
    exit 0
fi
