import { joinURL } from "ufo";
import { createMatcherFromExport } from "radix3";
import { defu } from "defu";
import { useAppConfig, useRuntimeConfig } from "#app";
import { appManifest as isAppManifestEnabled } from "#build/nuxt.config.mjs";
let manifest;
let matcher;
function fetchManifest() {
  if (!isAppManifestEnabled) {
    throw new Error("[nuxt] app manifest should be enabled with `experimental.appManifest`");
  }
  const config = useRuntimeConfig();
  const buildId = useAppConfig().nuxt?.buildId;
  manifest = $fetch(joinURL(config.app.cdnURL || config.app.baseURL, config.app.buildAssetsDir, `builds/meta/${buildId}.json`));
  manifest.then((m) => {
    matcher = createMatcherFromExport(m.matcher);
  });
  return manifest;
}
export function getAppManifest() {
  if (!isAppManifestEnabled) {
    throw new Error("[nuxt] app manifest should be enabled with `experimental.appManifest`");
  }
  return manifest || fetchManifest();
}
export async function getRouteRules(url) {
  await getAppManifest();
  return defu({}, ...matcher.matchAll(url).reverse());
}
