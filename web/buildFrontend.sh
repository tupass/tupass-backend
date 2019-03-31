#!/bin/bash
DIR="$(cd "$(dirname "$0")" && pwd)"

if ! [ -x "$(command -v npm)" ]; then
    echo "Unable to find npm, required to build frontend"
    exit 1
fi

if ! [ -x "$(command -v git)" ]; then
    echo "Unable to find git, required to clone frontend"
    exit 1
fi

rm -rf frontend
git clone https://github.com/tupass/tupass-frontend.git
cd tupass-frontend
npm install
npm run build-bundling
npm run build-bundling-de
mv dist/tupass-frontend/ ../frontend
cd $DIR
rm -rf tupass-frontend
rm frontend/de/assets/legal_notes.html frontend/de/assets/privacy_policy.html
rm frontend/en/assets/legal_notes.html frontend/en/assets/privacy_policy.html
