import htmx from "htmx.org";
import { initEpisodeViewer } from "./features/episode-viewer";
import { initResultsInteractions } from "./features/results-interactions";
import { initSearchPanel } from "./features/search-panel";
import { initSidebar } from "./features/sidebar";
import { initTheme } from "./features/theme";

(window as Window & { htmx?: unknown }).htmx = htmx;

document.addEventListener("DOMContentLoaded", (): void => {
    const root: HTMLElement = document.documentElement;
    const maxQueryChars: number = 500;
    const copyDebounceMs: number = 350;
    const copyFlashMs: number = 1200;

    initTheme(root, "pone-theme");
    initSidebar({
        minSidebarWidth: 320,
        maxSidebarWidth: 600,
    });

    initSearchPanel({
        maxQueryChars,
    });

    initResultsInteractions({
        copyDebounceMs,
        copyFlashMs,
        pinnedCardClass: "result-card--pinned",
    });

    initEpisodeViewer({
        episodeLineHighlightClass: "episode-line--highlight",
    });
});
