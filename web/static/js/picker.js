class Picker extends HTMLElement {
  connectedCallback() {
    this.listeners?.abort();
    this.listeners = new AbortController();
    const { signal } = this.listeners;
    const itemKey = this.getAttribute("item-key");
    const multiple = this.hasAttribute("multiple");
    const selectedVariant =
      this.getAttribute("item-selected-variant") || "badge";
    const details = this.querySelector("details");
    const search = this.querySelector("[data-picker-search]");
    const selectedItems = this.querySelector("[data-picker-selected]");
    const placeholder = this.querySelector("[data-picker-placeholder]");
    const empty = this.querySelector("[data-picker-empty]");
    const options = Array.from(
      this.querySelectorAll("[data-picker-option]"),
      (element) => {
        const inp = element.querySelector("input");
        return {
          element,
          input: inp,
          item: inp.nextElementSibling,
        };
      },
    );
    const badges = new Map(
      Array.from(selectedItems.children, (badge) => [
        badge.getAttribute(itemKey),
        badge,
      ]),
    );

    const filter = () => {
      const query = search.value.trim().toLocaleLowerCase();
      for (const { element, item } of options) {
        element.hidden = !item
          .getAttribute("name")
          .toLocaleLowerCase()
          .includes(query);
      }
      empty.hidden = options.some(({ element }) => !element.hidden);
    };

    const updateSelection = () => {
      for (const { input, item } of options) {
        if (input.checked && !badges.has(input.value)) {
          const badge = item.cloneNode();
          badge.setAttribute("variant", selectedVariant);
          badges.set(input.value, badge);
          selectedItems.append(badge);
        } else if (!input.checked && badges.has(input.value)) {
          badges.get(input.value).remove();
          badges.delete(input.value);
        }
      }
      selectedItems.hidden = badges.size === 0;
      placeholder.hidden = badges.size > 0;
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
        if (!options.some(({ input }) => input === event.target)) return;
        updateSelection();
        if (!multiple) details.open = false;
      },
      { signal },
    );
    this.addEventListener(
      "click",
      (event) => {
        if (
          !multiple &&
          options.some(({ input }) => input === event.target) &&
          badges.has(event.target.value)
        ) {
          details.open = false;
        }
      },
      { signal },
    );
    this.addEventListener(
      "remove",
      (event) => {
        const option = options.find(
          ({ input }) => input.value === event.detail.id,
        );
        if (!option) return;
        option.input.checked = false;
        option.input.dispatchEvent(new Event("change", { bubbles: true }));
      },
      { signal },
    );
    this.addEventListener(
      "keydown",
      (event) => {
        if (event.key === "Escape" && details.open) {
          event.preventDefault();
          details.open = false;
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

customElements.define("options-picker", Picker);
