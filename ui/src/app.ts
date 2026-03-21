import htmx from "htmx.org";

(window as Window & { htmx?: unknown }).htmx = htmx;

document.addEventListener("DOMContentLoaded", (): void => {
    const $ = (id: string): HTMLElement | null => document.getElementById(id);
    const root: HTMLElement = document.documentElement;
    const maxQueryChars: number = 500;
    const copyDebounceMs: number = 350;
    const copyFlashMs: number = 1200;
    const minSidebarWidth: number = 320;
    const maxSidebarWidth: number = 600;
    const themeStorageKey: string = "pone-theme";

    const themeToggleBtn: HTMLElement | null = $("theme-toggle");
    const themeToggleIcon: HTMLElement | null = $("theme-toggle-icon");

    const getPreferredTheme = (): "light" | "dark" => {
        const saved: string | null = window.localStorage.getItem(themeStorageKey);
        if (saved === "dark" || saved === "light") return saved;
        return "light";
    };

    const updateThemeToggleUI = (theme: "light" | "dark"): void => {
        if (themeToggleIcon) themeToggleIcon.textContent = theme === "dark" ? "☀" : "🌙";
        if (themeToggleBtn) {
            const nextModeLabel: string = theme === "dark" ? "切换到日间模式" : "切换到夜间模式";
            themeToggleBtn.setAttribute("aria-label", nextModeLabel);
            themeToggleBtn.setAttribute("title", nextModeLabel);
        }
    };

    const applyTheme = (theme: "light" | "dark"): void => {
        root.setAttribute("data-theme", theme);
        window.localStorage.setItem(themeStorageKey, theme);
        updateThemeToggleUI(theme);
    };

    applyTheme(getPreferredTheme());

    if (themeToggleBtn) {
        themeToggleBtn.addEventListener("click", (): void => {
            const current: "light" | "dark" = (root.getAttribute("data-theme") === "dark") ? "dark" : "light";
            applyTheme(current === "dark" ? "light" : "dark");
        });
    }

    const closeEpisodeViewer = (): void => {
        const viewer: HTMLElement | null = $("episode-viewer");
        if (!viewer) return;
        viewer.innerHTML = "";
        document.body.style.overflow = "";
    };

    const sidebar: HTMLElement | null = $("sidebar");
    const handle: HTMLElement | null = $("drag-handle");
    if (sidebar && handle) {
        let isResizing: boolean = false;
        handle.onmousedown = (e: MouseEvent): void => {
            isResizing = true;
            document.body.style.cursor = "col-resize";
            e.preventDefault();
        };
        document.onmousemove = (e: MouseEvent): void => {
            if (!isResizing) return;
            sidebar.style.width = `${Math.min(Math.max(e.clientX, minSidebarWidth), maxSidebarWidth)}px`;
        };
        document.onmouseup = (): void => {
            isResizing = false;
            document.body.style.cursor = "";
        };
    }

    const toggleBtn: HTMLElement | null = $("toggle-search");
    const searchContent: HTMLElement | null = $("search-content");
    const sidebarBrand: HTMLElement | null = $("sidebar-brand");
    const sidebarMeta: HTMLElement | null = $("sidebar-meta");
    let isMobileCollapsed: boolean = false;
    const mobileSections: (HTMLElement | null)[] = [searchContent, sidebarBrand, sidebarMeta];

    const setMobileCollapsed = (collapsed: boolean): void => {
        isMobileCollapsed = collapsed;
        mobileSections.forEach((el): void => {
            if (!el) return;
            el.classList.toggle("hidden", collapsed);
            (el as HTMLElement).style.display = collapsed ? "none" : "";
        });
    };

    if (toggleBtn && searchContent) {
        toggleBtn.addEventListener("click", () => {
            setMobileCollapsed(!isMobileCollapsed);
        });
    }

    const slider: HTMLInputElement | null = $("topk-slider") as HTMLInputElement | null;
    const realInput: HTMLInputElement | null = $("topk-real") as HTMLInputElement | null;
    const progress: HTMLElement | null = $("slider-progress") as HTMLElement | null;
    const display: HTMLElement | null = $("topk-display") as HTMLElement | null;
    const steps: number[] = [10, 25, 50, 100];

    const updateSlider = (val: string): void => {
        const idx: number = Number.parseInt(val, 10);
        const actualValue: number = steps[idx];
        if (realInput) realInput.value = String(actualValue);
        if (display) display.textContent = String(actualValue);
        if (progress) progress.style.width = `${(idx / 3) * 100}%`;
    };

    if (slider) {
        updateSlider(slider.value);
        slider.oninput = (e: Event): void => updateSlider((e.target as HTMLInputElement).value);
    }

    const searchInput: HTMLTextAreaElement | null = $("search-input") as HTMLTextAreaElement | null;
    const clearSearch: HTMLElement | null = $("clear-search");
    const charCount: HTMLElement | null = $("search-char-count");

    const toChars = (text: string | null | undefined): string[] => Array.from(text ?? "");

    const fallbackCopyText = (text: string): void => {
        const temp: HTMLTextAreaElement = document.createElement("textarea");
        temp.value = text;
        temp.style.position = "fixed";
        temp.style.top = "0";
        temp.style.left = "-9999px";
        temp.style.opacity = "0";
        document.body.appendChild(temp);
        temp.focus();
        temp.select();
        temp.setSelectionRange(0, temp.value.length);

        const ok: boolean = document.execCommand("copy");
        document.body.removeChild(temp);
        if (!ok) throw new Error("copy failed");
    };

    const copyText = async (text: string): Promise<void> => {
        if (navigator.clipboard?.writeText) {
            try {
                await navigator.clipboard.writeText(text);
                return;
            } catch {
                // Fallback for environments where clipboard API is blocked.
            }
        }
        fallbackCopyText(text);
    };

    const flashCopySuccess = (button: Element): void => {
        const copyIcon: Element | null = button.querySelector(".copy-icon");
        const checkIcon: Element | null = button.querySelector(".check-icon");
        if (!copyIcon || !checkIcon) return;
        copyIcon.classList.add("hidden");
        checkIcon.classList.remove("hidden");

        setTimeout((): void => {
            copyIcon.classList.remove("hidden");
            checkIcon.classList.add("hidden");
        }, copyFlashMs);
    };

    const updateSearchCounter = (): void => {
        if (!searchInput) return;

        let chars: string[] = toChars(searchInput.value);
        if (chars.length > maxQueryChars) {
            chars = chars.slice(0, maxQueryChars);
            searchInput.value = chars.join("");
        }

        if (charCount) {
            const current: number = chars.length;
            charCount.textContent = `${current}/${maxQueryChars}`;
            charCount.classList.toggle("text-brand-accent", current >= maxQueryChars);
            charCount.classList.toggle("font-semibold", current >= maxQueryChars);
        }
    };

    if (searchInput) {
        updateSearchCounter();
        searchInput.addEventListener("input", updateSearchCounter);
    }

    if (clearSearch && searchInput) {
        clearSearch.addEventListener("click", (): void => {
            searchInput.value = "";
            searchInput.focus();
            updateSearchCounter();
        });
    }

    document.addEventListener("keydown", (e: KeyboardEvent): void => {
        if (e.key !== "Escape") return;
        if ($("episode-viewer")?.children.length) {
            closeEpisodeViewer();
            return;
        }

        const active: HTMLTextAreaElement | null = document.activeElement as HTMLTextAreaElement | null;
        if (active?.id === "search-input") {
            active.value = "";
            updateSearchCounter();
        }
    });

    let lastCopyAt: number = 0;

    const handleGlobalTap = (e: Event): void => {
        const target: Element | null = e.target as Element | null;
        const copyButton: Element | null = target?.closest("[data-copy-line='true']") ?? null;
        if (copyButton) {
            const now: number = Date.now();
            if (now - lastCopyAt < copyDebounceMs) return;
            lastCopyAt = now;

            e.preventDefault();

            const lineText: string = copyButton.getAttribute("data-copy-text") ?? "";
            const character: string = copyButton.getAttribute("data-copy-character") ?? "Unknown";
            const seasonCode: string = copyButton.getAttribute("data-copy-seasoncode") ?? "unknown";
            const payload: string = `${lineText}\n  -- ${character}, ${seasonCode}`;

            copyText(payload)
                .then((): void => flashCopySuccess(copyButton))
                .catch((): void => {
                    window.prompt("复制失败，请手动复制以下内容：", payload);
                });
            return;
        }

        if (target?.closest("[data-episode-close='true']")) {
            closeEpisodeViewer();
        }
    };

    document.addEventListener("click", handleGlobalTap);
    document.addEventListener("touchend", handleGlobalTap, { passive: false });

    document.body.addEventListener("htmx:afterSwap", (e: Event): void => {
        const target: HTMLElement | null = e.target as HTMLElement | null;
        if (!target || target.id !== "episode-viewer") return;

        document.body.style.overflow = "hidden";
        requestAnimationFrame((): void => {
            const body: HTMLElement | null = target.querySelector(".episode-modal__body") as HTMLElement | null;
            const highlighted: HTMLElement | null = target.querySelector("#episode-highlight") as HTMLElement | null;
            if (highlighted && body) {
                const offsetTop: number = highlighted.offsetTop - body.offsetTop;
                const containerHeight: number = body.clientHeight;
                const elementHeight: number = highlighted.clientHeight;
                body.scrollTop = offsetTop - (containerHeight - elementHeight) / 2;
            }
        });
    });
});
