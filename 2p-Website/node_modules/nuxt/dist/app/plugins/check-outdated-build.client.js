import { joinURL } from "ufo";
import { defineNuxtPlugin, getAppManifest, onNuxtReady, useRuntimeConfig } from "#app";
export default defineNuxtPlugin((nuxtApp) => {
  if (import.meta.test) {
    return;
  }
  let timeout;
  const config = useRuntimeConfig();
  async function getLatestManifest() {
    const currentManifest = await getAppManifest();
    if (timeout) {
      clearTimeout(timeout);
    }
    timeout = setTimeout(getLatestManifest, 1e3 * 60 * 60);
    const meta = await $fetch(joinURL(config.app.cdnURL || config.app.baseURL, config.app.buildAssetsDir, "builds/latest.json"));
    if (meta.id !== currentManifest.id) {
      nuxtApp.hooks.callHook("app:manifest:update", meta);
    }
  }
  onNuxtReady(() => {
    timeout = setTimeout(getLatestManifest, 1e3 * 60 * 60);
  });
});
