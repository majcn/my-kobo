docker run --mount type=bind,source="$(pwd)/koboSync",target=/project --workdir="/project" majcn/e-reader-toolchain kobo
mv koboSync/koboSync sdcard/mnt/onboard/.adds/majcn/koboSync

docker run --mount type=bind,source="$(pwd)/translate",target=/project --workdir="/project" majcn/e-reader-toolchain kobo
mv translate/translate sdcard/mnt/onboard/.adds/majcn/translate
