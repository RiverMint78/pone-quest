import htmx from "htmx.org";
import { initBackToTopButton } from "./features/back-to-top";
import { initEpisodeViewer } from "./features/episode-viewer";
import { initResultsInteractions } from "./features/results-interactions";
import { initSearchPanel } from "./features/search-panel";
import { initSidebar } from "./features/sidebar";
import { initTheme } from "./features/theme";

(window as Window & { htmx?: unknown }).htmx = htmx;

document.addEventListener("DOMContentLoaded", (): void => {
    initTheme();
    initSidebar();
    initSearchPanel();
    initResultsInteractions();
    initEpisodeViewer();
    initBackToTopButton();
});
