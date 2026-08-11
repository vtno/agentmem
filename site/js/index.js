/**
 * agentmem site — code copy buttons + theme toggle.
 * Copy enhance runs on full load and after htmx #main swaps.
 * Theme control lives outside #main so it survives swaps.
 */
(function () {
  "use strict";

  var ATTR = "data-copy-ready";
  var THEME_KEY = "agentmem-theme";

  var ICON_COPY =
    '<svg class="copy-btn__icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" aria-hidden="true">' +
    '<rect x="9" y="9" width="13" height="13" rx="0"/>' +
    '<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>' +
    "</svg>";

  var ICON_CHECK =
    '<svg class="copy-btn__icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" aria-hidden="true">' +
    '<path d="M20 6L9 17l-5-5"/>' +
    "</svg>";

  /* Sun — shown in dark mode (action: switch to light) */
  var ICON_SUN =
    '<svg class="theme-toggle__icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" aria-hidden="true">' +
    '<circle cx="12" cy="12" r="4"/>' +
    '<path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"/>' +
    "</svg>";

  /* Moon — shown in light mode (action: switch to dark) */
  var ICON_MOON =
    '<svg class="theme-toggle__icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" aria-hidden="true">' +
    '<path d="M21 14.5A8.5 8.5 0 0 1 9.5 3 7 7 0 1 0 21 14.5z"/>' +
    "</svg>";

  function textOf(block) {
    var code = block.querySelector("code");
    var raw = code ? code.textContent : block.textContent;
    return (raw || "").replace(/\n$/, "");
  }

  function setIcon(btn, kind) {
    btn.innerHTML = kind === "ok" ? ICON_CHECK : ICON_COPY;
    btn.setAttribute("aria-label", kind === "ok" ? "Copied" : "Copy code");
  }

  function popMascot() {
    var img = document.createElement("img");
    img.className = "copy-pop";
    img.alt = "";
    /* Inverted theme → file mapping (matches the brand mascot
     * convention in layout.tmpl + syncBrandMascot): dark theme
     * shows the LIGHT-coloured robot, light theme shows the
     * DARK-coloured one, so the icon reads against its bg. */
    var name = currentTheme() === "dark" ? "light" : "dark";
    img.src = "/assets/icons/agentmem-mascot-thumb-up-" + name + ".png";
    img.addEventListener("animationend", function () {
      if (img.parentNode) img.parentNode.removeChild(img);
    });
    document.body.appendChild(img);
  }

  function copyText(text) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
      return navigator.clipboard.writeText(text);
    }
    return new Promise(function (resolve, reject) {
      var ta = document.createElement("textarea");
      ta.value = text;
      ta.setAttribute("readonly", "");
      ta.style.position = "fixed";
      ta.style.left = "-9999px";
      document.body.appendChild(ta);
      ta.select();
      try {
        if (!document.execCommand("copy")) {
          reject(new Error("copy failed"));
        } else {
          resolve();
        }
      } catch (err) {
        reject(err);
      } finally {
        document.body.removeChild(ta);
      }
    });
  }

  function enhanceBlock(block) {
    if (block.getAttribute(ATTR) === "1") return;
    // Skip if already wrapped (e.g. partial re-enhance)
    if (block.parentElement && block.parentElement.classList.contains("code-block")) {
      block.setAttribute(ATTR, "1");
      return;
    }
    block.setAttribute(ATTR, "1");

    var wrap = document.createElement("div");
    wrap.className = "code-block";
    block.parentNode.insertBefore(wrap, block);
    wrap.appendChild(block);

    var btn = document.createElement("button");
    btn.type = "button";
    btn.className = "copy-btn";
    setIcon(btn, "copy");

    btn.addEventListener("click", function () {
      var text = textOf(block);
      copyText(text)
        .then(function () {
          setIcon(btn, "ok");
          popMascot();
          window.setTimeout(function () {
            setIcon(btn, "copy");
          }, 1400);
        })
        .catch(function () {
          btn.setAttribute("aria-label", "Copy failed");
          window.setTimeout(function () {
            setIcon(btn, "copy");
          }, 1400);
        });
    });

    wrap.appendChild(btn);
  }

  function enhanceAll(root) {
    var scope = root && root.querySelectorAll ? root : document;
    var blocks = scope.querySelectorAll("pre");
    for (var i = 0; i < blocks.length; i++) {
      enhanceBlock(blocks[i]);
    }
  }

  function hasExplicitTheme() {
    var t = document.documentElement.getAttribute("data-theme");
    return t === "light" || t === "dark";
  }

  function systemTheme() {
    try {
      if (window.matchMedia("(prefers-color-scheme: dark)").matches) return "dark";
    } catch (e) {}
    return "light";
  }

  function currentTheme() {
    if (hasExplicitTheme()) {
      return document.documentElement.getAttribute("data-theme");
    }
    return systemTheme();
  }

  function applyTheme(theme) {
    document.documentElement.setAttribute("data-theme", theme);
    try {
      localStorage.setItem(THEME_KEY, theme);
    } catch (e) {}
    syncThemeButton(theme);
    syncBrandMascot(theme);
    syncLedeMascot(theme);
  }

  function syncBrandMascot(theme) {
    var img = document.getElementById("brand-mascot");
    if (!img) return;
    /* File suffix names the mascot's fill colour, not the theme. The
     * `-light` mascot is a LIGHT-coloured robot (visible on dark bg),
     * the `-dark` mascot is a DARK-coloured one (visible on light bg).
     * Invert the mapping so each mascot reads against its bg. */
    var name = theme === "dark" ? "light" : "dark";
    var next = "/assets/icons/agentmem-mascot-normal-" + name + ".png";
    /* Only swap if it actually changed — avoids reloading the image. */
    if (img.getAttribute("src") !== next) img.setAttribute("src", next);
  }

  function syncLedeMascot(theme) {
    var img = document.getElementById("lede-mascot");
    if (!img) return;
    /* Same inverted theme → file mapping as syncBrandMascot. The lede
     * uses the heart variant. */
    var name = theme === "dark" ? "light" : "dark";
    var next = "/assets/icons/agentmem-mascot-heart-" + name + ".png";
    if (img.getAttribute("src") !== next) img.setAttribute("src", next);
  }

  function syncThemeButton(theme) {
    var btn = document.getElementById("theme-toggle");
    if (!btn) return;
    var isDark = theme === "dark";
    btn.innerHTML = isDark ? ICON_SUN : ICON_MOON;
    var label = isDark
      ? btn.getAttribute("data-label-light") || "Switch to light mode"
      : btn.getAttribute("data-label-dark") || "Switch to dark mode";
    btn.setAttribute("aria-label", label);
    btn.setAttribute("title", label);
  }

  function syncThemeVisuals() {
    var theme = currentTheme();
    syncBrandMascot(theme);
    syncLedeMascot(theme);
    syncThemeButton(theme);
  }

  function initThemeToggle() {
    var btn = document.getElementById("theme-toggle");
    if (!btn || btn.getAttribute("data-theme-ready") === "1") return;
    btn.setAttribute("data-theme-ready", "1");
    syncThemeVisuals();
    btn.addEventListener("click", function () {
      applyTheme(currentTheme() === "dark" ? "light" : "dark");
    });
    /* While following system, keep the icon in sync if OS theme changes. */
    try {
      var mq = window.matchMedia("(prefers-color-scheme: dark)");
      var onChange = function () {
        if (!hasExplicitTheme()) {
          var sys = systemTheme();
          syncThemeButton(sys);
          syncBrandMascot(sys);
          syncLedeMascot(sys);
        }
      };
      if (mq.addEventListener) mq.addEventListener("change", onChange);
      else if (mq.addListener) mq.addListener(onChange);
    } catch (e) {}
  }

  function currentLang() {
    var btn = document.getElementById("lang-toggle");
    var lang = btn && btn.getAttribute("data-lang");
    lang = (lang || "en").toLowerCase();
    return lang === "th" ? "th" : "en";
  }

  function syncLangButton() {
    var btn = document.getElementById("lang-toggle");
    if (!btn) return;
    var cur = currentLang();
    var next = cur === "th" ? "en" : "th";
    var code =
      next === "th"
        ? btn.getAttribute("data-code-th") || "TH"
        : btn.getAttribute("data-code-en") || "EN";
    var label =
      next === "th"
        ? btn.getAttribute("data-label-th") || "สลับเป็นภาษาไทย"
        : btn.getAttribute("data-label-en") || "Switch to English";
    btn.textContent = code;
    btn.setAttribute("aria-label", label);
    btn.setAttribute("title", label);
  }

  /* Set (or replace) lang= via URLSearchParams — never append a second copy. */
  function applyLangParam(url, lang) {
    url.searchParams.set("lang", lang);
    return url;
  }

  /* Keep internal navigations on the active locale for htmx boost,
   * preload, and history URLs. Skip external, hash, asset, and .md links. */
  function withLang(href, lang) {
    if (!href || href.charAt(0) === "#" || href.indexOf("mailto:") === 0) {
      return href;
    }
    try {
      var u = new URL(href, window.location.origin);
      if (u.origin !== window.location.origin) return href;
      if (
        u.pathname.indexOf("/css/") === 0 ||
        u.pathname.indexOf("/js/") === 0 ||
        u.pathname.indexOf("/assets/") === 0
      ) {
        return href;
      }
      if (/\.md$/i.test(u.pathname)) return href;
      applyLangParam(u, lang);
      return u.pathname + u.search + u.hash;
    } catch (e) {
      return href;
    }
  }

  function stampLangLinks(root) {
    var lang = currentLang();
    var scope = root && root.querySelectorAll ? root : document;
    var links = scope.querySelectorAll("a[href]");
    for (var i = 0; i < links.length; i++) {
      var a = links[i];
      if (a.getAttribute("hx-boost") === "false") continue;
      var href = a.getAttribute("href");
      var next = withLang(href, lang);
      if (next !== href) a.setAttribute("href", next);
    }
  }

  function persistLangCookie(lang) {
    try {
      document.cookie =
        "lang=" + lang + ";path=/;max-age=31536000;SameSite=Lax";
    } catch (e) {}
  }

  function initLangToggle() {
    var btn = document.getElementById("lang-toggle");
    if (!btn || btn.getAttribute("data-lang-ready") === "1") return;
    btn.setAttribute("data-lang-ready", "1");
    /* Cookie is set client-side so HTML with ?lang= stays free of
     * Set-Cookie and Cloudflare can cache EN/TH as separate URLs. */
    persistLangCookie(currentLang());
    syncLangButton();
    stampLangLinks(document);
    btn.addEventListener("click", function () {
      var next = currentLang() === "th" ? "en" : "th";
      persistLangCookie(next);
      /* Full navigation: chrome (lang toggle) is outside the htmx
       * #content partial, so locale must re-render the whole page. */
      var url = applyLangParam(new URL(window.location.href), next);
      window.location.assign(url.pathname + url.search + url.hash);
    });
  }

  function onReady() {
    initLangToggle();
    initThemeToggle();
    enhanceAll(document);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", onReady);
  } else {
    onReady();
  }

  /* Ensure the request path has exactly one lang= (set/replace). Do not
   * also put lang on detail.parameters — htmx would append a second copy
   * when the href already includes ?lang=. */
  document.body.addEventListener("htmx:configRequest", function (evt) {
    if (!evt.detail || !evt.detail.path) return;
    try {
      var u = new URL(evt.detail.path, window.location.origin);
      if (u.origin !== window.location.origin) return;
      applyLangParam(u, currentLang());
      evt.detail.path = u.pathname + u.search + u.hash;
      if (evt.detail.parameters && "lang" in evt.detail.parameters) {
        delete evt.detail.parameters.lang;
      }
    } catch (e) {}
  });

  document.body.addEventListener("htmx:afterSettle", function (evt) {
    var target = evt.detail && evt.detail.target;
    /* The server-rendered partial does not know the client's theme. */
    syncThemeVisuals();
    stampLangLinks(target || document);
    enhanceAll(target || document);
  });
})();
