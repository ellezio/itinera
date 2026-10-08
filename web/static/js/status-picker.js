class ResourceStatusPicker extends HTMLElement {
  connectedCallback() {
    this.listeners?.abort();
    this.listeners = new AbortController();
    const { signal } = this.listeners;
    const details = this.querySelector("details");
    const search = this.querySelector("[data-status-search]");
    const options = Array.from(this.querySelectorAll("[data-status-option]"));
    const selectedStatus = this.querySelector("[data-status-selected]");

    const filter = () => {
      const query = search.value.trim().toLocaleLowerCase();
      for (const option of options) {
        const name = option
          .querySelector("resource-status")
          .getAttribute("name");
        option.hidden = !name.toLocaleLowerCase().includes(query);
      }
      this.querySelector("[data-status-empty]").hidden = options.some(
        (option) => !option.hidden,
      );
    };

    const updateSelection = () => {
      const option = options.find(
        (option) => option.querySelector("input").checked,
      );
      if (option) {
        const status = option.querySelector("resource-status");
        let badge = selectedStatus.querySelector("resource-status");
        if (!badge) {
          badge = document.createElement("resource-status");
          badge.setAttribute("variant", "badge");
          selectedStatus.append(badge);
        }
        badge.status = {
          id: status.getAttribute("status-id"),
          name: status.getAttribute("name"),
          color: status.getAttribute("color"),
        };
      } else {
        selectedStatus.replaceChildren();
      }
      selectedStatus.hidden = !option;
      this.querySelector("[data-status-placeholder]").hidden = Boolean(option);
    };

    const close = () => {
      details.open = false;
    };

    details.addEventListener(
      "toggle",
      () => {
        if (details.open) {
          search.focus();
        } else {
          search.value = "";
          filter();
        }
      },
      { signal },
    );
    search.addEventListener("input", filter, { signal });
    this.addEventListener(
      "change",
      (event) => {
        if (!event.target.matches('input[type="radio"][name="status"]')) return;
        updateSelection();
        close();
      },
      { signal },
    );
    this.addEventListener(
      "click",
      (event) => {
        if (
          event.target.matches('input[type="radio"][name="status"]') &&
          selectedStatus.firstElementChild?.getAttribute("status-id") ===
            event.target.value
        ) {
          close();
        }
      },
      { signal },
    );
    this.addEventListener(
      "keydown",
      (event) => {
        if (event.key === "Escape" && details.open) {
          event.preventDefault();
          close();
        } else if (event.target === search && event.key === "Enter") {
          event.preventDefault();
        }
      },
      { signal },
    );
    document.addEventListener(
      "click",
      (event) => {
        if (!this.contains(event.target)) details.open = false;
      },
      { signal },
    );
    this.addEventListener(
      "focusout",
      (event) => {
        if (event.relatedTarget && !this.contains(event.relatedTarget)) {
          details.open = false;
        }
      },
      { signal },
    );

    filter();
    updateSelection();
  }

  disconnectedCallback() {
    this.listeners.abort();
  }
}

customElements.define("resource-status-picker", ResourceStatusPicker);
