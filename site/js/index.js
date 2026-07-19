/**
 * agentmem site — enhance code blocks with a copy icon button.
 * Works on full load and after htmx #main swaps.
 */
(function () {
  "use strict";

  var ATTR = "data-copy-ready";

  var ICON_COPY =
    '<svg class="copy-btn__icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" aria-hidden="true">' +
    '<rect x="9" y="9" width="13" height="13" rx="0"/>' +
    '<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>' +
    "</svg>";

  var ICON_CHECK =
    '<svg class="copy-btn__icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="square" stroke-linejoin="miter" aria-hidden="true">' +
    '<path d="M20 6L9 17l-5-5"/>' +
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

  function onReady() {
    enhanceAll(document);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", onReady);
  } else {
    onReady();
  }

  document.body.addEventListener("htmx:afterSettle", function (evt) {
    var target = evt.detail && evt.detail.target;
    enhanceAll(target || document);
  });
})();
