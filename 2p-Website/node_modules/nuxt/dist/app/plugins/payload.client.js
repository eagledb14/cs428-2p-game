import { parseURL } from "ufo";
import { defineNuxtPlugin } from "#app/nuxt";
import { loadPayload } from "#app/composables/payload";
import { onNuxtReady } from "#app/composables/ready";
import { useRouter } from "#app/composables/router";
import { getAppManifest } from "#app/composables/manifest";
import { appManifest as isAppManifestEnabled } from "#build/nuxt.config.mjs";
export default defineNuxtPlugin({
  name: "nuxt:payload",
  setup(nuxtApp) {
    if (process.dev) {
      return;
    }
    useRouter().beforeResolve(async (to, from) => {
      if (to.path === from.path) {
        return;
      }
      const payload = await loadPayload(to.path);
      if (!payload) {
        return;
      }
      Object.assign(nuxtApp.static.data, payload.data);
    });
    onNuxtReady(() => {
      nuxtApp.hooks.hook("link:prefetch", async (url) => {
        if (!parseURL(url).protocol) {
          await loadPayload(url);
        }
      });
      if (isAppManifestEnabled && navigator.connection?.effectiveType !== "slow-2g") {
        setTimeout(getAppManifest, 1e3);
      }
    });
  }
});
