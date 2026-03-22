interface EpisodeViewerOptions {
    episodeLineHighlightClass: string;
}

interface EpisodeViewerController {
    closeEpisodeViewer: () => void;
    hasOpenViewer: () => boolean;
}

export function initEpisodeViewer(options: EpisodeViewerOptions): EpisodeViewerController {
    const episodeViewer: HTMLElement | null = document.getElementById("episode-viewer");
    let bodyPrevOverflow: string = "";
    let bodyPrevPaddingRight: string = "";

    const lockBodyScroll = (): void => {
        const scrollbarWidth: number = window.innerWidth - document.documentElement.clientWidth;
        bodyPrevOverflow = document.body.style.overflow;
        bodyPrevPaddingRight = document.body.style.paddingRight;
        document.body.style.overflow = "hidden";
        if (scrollbarWidth > 0) {
            document.body.style.paddingRight = `${scrollbarWidth}px`;
        }
    };

    const unlockBodyScroll = (): void => {
        document.body.style.overflow = bodyPrevOverflow;
        document.body.style.paddingRight = bodyPrevPaddingRight;
    };

    const closeEpisodeViewer = (): void => {
        if (!episodeViewer) return;
        episodeViewer.innerHTML = "";
        unlockBodyScroll();
    };

    const hasOpenViewer = (): boolean => {
        return (episodeViewer?.children.length ?? 0) > 0;
    };

    const syncEpisodeHitMarkers = (episodeBody: HTMLElement): void => {
        if (!episodeViewer) return;

        const hitMarkers: HTMLElement | null = episodeViewer.querySelector("#episode-hit-markers") as HTMLElement | null;
        if (!hitMarkers) return;

        hitMarkers.innerHTML = "";
        const highlightedLines: NodeListOf<HTMLElement> = episodeBody.querySelectorAll(`.episode-line.${options.episodeLineHighlightClass}`);
        const maxScrollable: number = Math.max(episodeBody.scrollHeight - episodeBody.clientHeight, 1);

        highlightedLines.forEach((line): void => {
            const lineCenterOffset: number = (line.offsetTop - episodeBody.offsetTop) + line.clientHeight / 2;
            const centeredScrollTop: number = lineCenterOffset - episodeBody.clientHeight / 2;
            const ratio: number = Math.min(Math.max(centeredScrollTop / maxScrollable, 0), 1);

            const marker: HTMLDivElement = document.createElement("div");
            marker.className = "absolute top-1/2 h-2.5 w-2.5 -translate-x-1/2 -translate-y-1/2 rounded-full bg-brand-accent border border-cream-bg";
            marker.style.left = `${Math.round(ratio * 100)}%`;
            hitMarkers.appendChild(marker);
        });
    };

    if (episodeViewer) {
        episodeViewer.addEventListener("click", (e: Event): void => {
            const target: Element | null = e.target as Element | null;

            if (target?.closest("[data-episode-close='true']")) {
                e.preventDefault();
                e.stopPropagation();
                closeEpisodeViewer();
                return;
            }

            const episodeLine: HTMLElement | null = target?.closest(".episode-line") as HTMLElement | null;
            if (!episodeLine) return;

            episodeLine.classList.toggle(options.episodeLineHighlightClass);
            const episodeBody: HTMLElement | null = episodeLine.closest(".episode-modal__body") as HTMLElement | null;
            if (episodeBody) {
                syncEpisodeHitMarkers(episodeBody);
            }
        });
    }

    document.body.addEventListener("htmx:afterSwap", (e: Event): void => {
        const target: HTMLElement | null = e.target as HTMLElement | null;
        if (!target || target.id !== "episode-viewer") return;

        lockBodyScroll();
        requestAnimationFrame((): void => {
            const body: HTMLElement | null = target.querySelector(".episode-modal__body") as HTMLElement | null;
            const highlighted: HTMLElement | null = target.querySelector("#episode-highlight") as HTMLElement | null;
            const progressBar: HTMLElement | null = target.querySelector("#episode-read-progress") as HTMLElement | null;

            const updateEpisodeReadProgress = (): void => {
                if (!body || !progressBar) return;
                const maxScrollable: number = Math.max(body.scrollHeight - body.clientHeight, 0);
                const ratio: number = maxScrollable === 0 ? 1 : Math.min(Math.max(body.scrollTop / maxScrollable, 0), 1);
                progressBar.style.width = `${Math.round(ratio * 100)}%`;
            };

            if (highlighted && body) {
                const offsetTop: number = highlighted.offsetTop - body.offsetTop;
                const containerHeight: number = body.clientHeight;
                const elementHeight: number = highlighted.clientHeight;
                body.scrollTop = offsetTop - (containerHeight - elementHeight) / 2;
            }

            if (body) {
                syncEpisodeHitMarkers(body);
            }
            updateEpisodeReadProgress();
            body?.addEventListener("scroll", updateEpisodeReadProgress, { passive: true });
        });
    });

    return {
        closeEpisodeViewer,
        hasOpenViewer,
    };
}
