const statusColors = new Set([
  "red",
  "green",
  "yellow",
  "blue",
  "purple",
  "aqua",
  "gray",
  "orange",
]);

class ResourceStatus extends HTMLElement {
  static observedAttributes = ["status-id", "name", "color", "variant"];

  connectedCallback() {
    if (!this.nameElement) {
      this.innerHTML = `
        <span class="status-marker" aria-hidden="true"></span>
        <span class="status-name"></span>
      `;
      this.nameElement = this.querySelector(".status-name");
      this.markerElement = this.querySelector(".status-marker");
    }
    this.render();
  }

  attributeChangedCallback() {
    if (this.nameElement) this.render();
  }

  get status() {
    return {
      id: this.getAttribute("status-id") || "",
      name: this.getAttribute("name") || "",
      color: this.getAttribute("color") || "gray",
    };
  }

  set status(metadata) {
    this.setAttribute("status-id", metadata.id ?? "");
    this.setAttribute("name", metadata.name ?? "");
    this.setAttribute("color", metadata.color ?? "gray");
  }

  render() {
    const { name, color } = this.status;
    const themeColor = statusColors.has(color) ? color : "gray";
    this.style.setProperty("--status-color", `var(--color-${themeColor}-1)`);
    this.style.setProperty(
      "--status-background",
      `var(--color-${themeColor}-0)`,
    );
    this.nameElement.textContent = name;
    this.markerElement.hidden =
      (this.getAttribute("variant") || "plain") !== "plain";
  }
}

customElements.define("resource-status", ResourceStatus);
