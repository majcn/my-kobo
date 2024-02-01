.PHONY: clean bin deploy

TARGET_RELEASE_PATH=target/arm-unknown-linux-musleabihf/release
DEPLOY_PATH=sdcard/mnt/onboard/.adds/majcn

clean:
	rm -rf target/

compile: $(TARGET_RELEASE_PATH)/kobo_github_sync $(TARGET_RELEASE_PATH)/translate_google $(TARGET_RELEASE_PATH)/translate_googlefree $(TARGET_RELEASE_PATH)/translate_termania

deploy: $(DEPLOY_PATH)/kobo_github_sync $(DEPLOY_PATH)/translate_google $(DEPLOY_PATH)/translate_googlefree $(DEPLOY_PATH)/translate_termania

$(TARGET_RELEASE_PATH)/kobo_github_sync:
$(TARGET_RELEASE_PATH)/translate_google:
$(TARGET_RELEASE_PATH)/translate_googlefree:
$(TARGET_RELEASE_PATH)/translate_termania:
	cross build --release --target=arm-unknown-linux-musleabihf

$(DEPLOY_PATH)/kobo_github_sync: $(TARGET_RELEASE_PATH)/kobo_github_sync
	cp $^ $@
$(DEPLOY_PATH)/translate_google: $(TARGET_RELEASE_PATH)/translate_google
	cp $^ $@
$(DEPLOY_PATH)/translate_googlefree: $(TARGET_RELEASE_PATH)/translate_googlefree
	cp $^ $@
$(DEPLOY_PATH)/translate_termania: $(TARGET_RELEASE_PATH)/translate_termania
	cp $^ $@

all: clean compile deploy
