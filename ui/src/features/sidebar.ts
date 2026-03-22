interface SidebarOptions {
    minSidebarWidth: number;
    maxSidebarWidth: number;
}

export function initSidebar(options: SidebarOptions): void {
    const sidebar: HTMLElement | null = document.getElementById("sidebar");
    const handle: HTMLElement | null = document.getElementById("drag-handle");

    if (sidebar && handle) {
        let isResizing: boolean = false;
        handle.onmousedown = (e: MouseEvent): void => {
            isResizing = true;
            document.body.style.cursor = "col-resize";
            e.preventDefault();
        };
        document.onmousemove = (e: MouseEvent): void => {
            if (!isResizing) return;
            sidebar.style.width = `${Math.min(Math.max(e.clientX, options.minSidebarWidth), options.maxSidebarWidth)}px`;
        };
        document.onmouseup = (): void => {
            isResizing = false;
            document.body.style.cursor = "";
        };
    }

    const toggleBtn: HTMLElement | null = document.getElementById("toggle-search");
    const searchContent: HTMLElement | null = document.getElementById("search-content");
    const sidebarBrand: HTMLElement | null = document.getElementById("sidebar-brand");
    const sidebarMeta: HTMLElement | null = document.getElementById("sidebar-meta");
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
}
