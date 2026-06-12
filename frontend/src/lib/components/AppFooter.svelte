<footer>
  <span class="brand title-lg">Blessed Bites</span>
  <nav aria-label="Footer">
    <a class="label-sm" href="/#about">Privacy Policy</a>
    <a class="label-sm" href="/#about">Terms of Service</a>
    <a class="label-sm" href="/#about">Contact Us</a>
  </nav>
  <span class="body-md copyright">© {new Date().getFullYear()} Blessed Bites. Premium Fast-Casual Dining.</span>
</footer>

<style>
  footer {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
    gap: var(--bb-space-md);
    width: 100%;
    padding: var(--bb-space-xl) var(--bb-margin-mobile);
    background: var(--md-sys-color-surface-container-highest);
    border-top: 1px solid var(--md-sys-color-outline-variant);
    margin-top: auto;
  }

  @media (min-width: 768px) {
    footer {
      flex-direction: row;
      padding: var(--bb-space-xl) var(--bb-margin-desktop);
    }
  }

  .brand {
    font-weight: 900;
    color: var(--md-sys-color-primary);
  }

  nav {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: var(--bb-space-lg);
  }

  nav a {
    color: var(--md-sys-color-on-surface-variant);
    text-decoration: none;
  }

  nav a:hover {
    color: var(--md-sys-color-primary);
    text-decoration: underline;
  }

  .copyright {
    color: var(--md-sys-color-on-surface);
    text-align: center;
  }
</style>
