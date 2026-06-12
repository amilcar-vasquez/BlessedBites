<script lang="ts">
  import { onMount } from 'svelte';
  import {
    deleteAdminUser,
    listAdminUsers,
    updateAdminUser,
    type AdminUser
  } from '$lib/api/admin';
  import { auth } from '$lib/stores/auth';
  import { showToast } from '$lib/stores/toast';
  import Skeleton from '$lib/components/Skeleton.svelte';

  let users = $state<AdminUser[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  let editing = $state<AdminUser | null>(null);
  let saving = $state(false);
  let form = $state({ email: '', full_name: '', phone_no: '', role: 'customer' });

  let deleting = $state<AdminUser | null>(null);

  async function refresh() {
    try {
      users = await listAdminUsers();
      error = null;
    } catch {
      error = 'Could not load users.';
    } finally {
      loading = false;
    }
  }

  onMount(refresh);

  function openEdit(user: AdminUser) {
    editing = user;
    form = {
      email: user.email,
      full_name: user.full_name,
      phone_no: user.phone_no,
      role: user.role
    };
  }

  async function save(e: Event) {
    e.preventDefault();
    if (!editing || saving) return;
    saving = true;
    try {
      await updateAdminUser(editing.id, {
        email: form.email.trim(),
        full_name: form.full_name.trim(),
        phone_no: form.phone_no.trim(),
        role: form.role
      });
      showToast(`${form.full_name} updated`, 'success');
      editing = null;
      await refresh();
    } catch {
      showToast('Could not update the user.', 'error');
    } finally {
      saving = false;
    }
  }

  async function confirmDelete() {
    if (!deleting) return;
    try {
      await deleteAdminUser(deleting.id);
      showToast(`${deleting.full_name} deleted`, 'success');
      deleting = null;
      await refresh();
    } catch {
      showToast('Could not delete the user.', 'error');
    }
  }

  function formatDate(iso: string): string {
    return new Date(iso).toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }
</script>

<svelte:head>
  <title>Users — Blessed Admin</title>
</svelte:head>

<header class="head">
  <h1 class="headline-lg">Users</h1>
  <p class="body-md muted">Manage customer and staff accounts.</p>
</header>

{#if error}
  <p class="bb-form-error" role="alert">{error}</p>
{/if}

{#if loading}
  <Skeleton width="100%" height="320px" radius="16px" />
{:else if users.length === 0}
  <div class="bb-card empty">
    <span class="material-symbols-outlined" aria-hidden="true">group_off</span>
    <p class="body-lg">No users found.</p>
  </div>
{:else}
  <div class="bb-card table-card">
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th class="label-sm">User</th>
            <th class="label-sm">Phone</th>
            <th class="label-sm">Role</th>
            <th class="label-sm">Joined</th>
            <th class="label-sm actions-col">Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each users as user (user.id)}
            <tr>
              <td>
                <div class="user-cell">
                  <span class="avatar label-sm">
                    {(user.full_name || '?')
                      .split(' ')
                      .map((p) => p[0])
                      .slice(0, 2)
                      .join('')
                      .toUpperCase()}
                  </span>
                  <div class="user-text">
                    <span class="title-md">{user.full_name}</span>
                    <span class="label-sm muted">{user.email}</span>
                  </div>
                </div>
              </td>
              <td class="body-md">{user.phone_no || '—'}</td>
              <td>
                <span class="chip label-sm" class:admin-chip={user.role === 'admin'}>
                  {user.role}
                </span>
              </td>
              <td class="body-md muted">{formatDate(user.created_at)}</td>
              <td>
                <div class="row-actions">
                  <button type="button" class="icon-btn" aria-label={`Edit ${user.full_name}`} onclick={() => openEdit(user)}>
                    <span class="material-symbols-outlined">edit</span>
                  </button>
                  <button
                    type="button"
                    class="icon-btn danger"
                    aria-label={`Delete ${user.full_name}`}
                    disabled={user.id === $auth.user?.id}
                    onclick={() => (deleting = user)}
                  >
                    <span class="material-symbols-outlined">delete</span>
                  </button>
                </div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
{/if}

<!-- Edit dialog -->
{#if editing}
  <div class="overlay" role="presentation" onclick={(e) => e.target === e.currentTarget && (editing = null)}>
    <div class="dialog bb-card" role="dialog" aria-modal="true" aria-label="Edit user">
      <h2 class="headline-md">Edit {editing.full_name}</h2>
      <form onsubmit={save}>
        <label class="bb-field">
          <span>Full name</span>
          <input type="text" required bind:value={form.full_name} />
        </label>
        <label class="bb-field">
          <span>Email</span>
          <input type="email" required bind:value={form.email} />
        </label>
        <label class="bb-field">
          <span>Phone number</span>
          <input type="tel" bind:value={form.phone_no} />
        </label>
        <label class="bb-field">
          <span>Role</span>
          <select bind:value={form.role}>
            <option value="customer">customer</option>
            <option value="admin">admin</option>
          </select>
        </label>
        <div class="dialog-actions">
          <button type="button" class="text-btn label-lg" onclick={() => (editing = null)}>Cancel</button>
          <button type="submit" class="bb-btn-primary" disabled={saving}>
            {saving ? 'Saving…' : 'Save changes'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Delete confirm -->
{#if deleting}
  <div class="overlay" role="presentation" onclick={(e) => e.target === e.currentTarget && (deleting = null)}>
    <div class="dialog bb-card" role="alertdialog" aria-modal="true" aria-label="Confirm deletion">
      <h2 class="headline-md">Delete {deleting.full_name}?</h2>
      <p class="body-md muted">
        This permanently removes the account for {deleting.email}. This action cannot be undone.
      </p>
      <div class="dialog-actions">
        <button type="button" class="text-btn label-lg" onclick={() => (deleting = null)}>Cancel</button>
        <button type="button" class="danger-btn label-lg" onclick={confirmDelete}>Delete</button>
      </div>
    </div>
  </div>
{/if}

<style>
  .head {
    margin-bottom: var(--bb-space-lg);
  }

  .head h1 {
    margin: 0;
    color: var(--md-sys-color-on-surface);
  }

  .head p {
    margin: var(--bb-space-xs) 0 0;
  }

  .muted {
    color: var(--md-sys-color-on-surface-variant);
  }

  .empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--bb-space-md);
    padding: var(--bb-space-xl);
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
  }

  .empty .material-symbols-outlined {
    font-size: 48px;
    color: var(--md-sys-color-outline);
  }

  .empty p {
    margin: 0;
  }

  .table-card {
    padding: var(--bb-space-md);
  }

  .table-wrap {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th {
    text-align: left;
    color: var(--md-sys-color-on-surface-variant);
    text-transform: uppercase;
    padding: 10px 12px;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  td {
    padding: 12px;
    color: var(--md-sys-color-on-surface);
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    vertical-align: middle;
  }

  tbody tr:nth-child(even) {
    background: var(--md-sys-color-surface-container-low);
  }

  tbody tr:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .actions-col {
    width: 96px;
  }

  .user-cell {
    display: flex;
    align-items: center;
    gap: var(--bb-space-md);
    min-width: 220px;
  }

  .avatar {
    width: 36px;
    height: 36px;
    border-radius: var(--bb-shape-full);
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    flex-shrink: 0;
  }

  .user-text {
    display: flex;
    flex-direction: column;
    gap: 2px;
    min-width: 0;
  }

  .chip {
    display: inline-block;
    padding: 4px 12px;
    border-radius: var(--bb-shape-full);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface-variant);
    text-transform: capitalize;
  }

  .chip.admin-chip {
    border-color: transparent;
    background: var(--md-sys-color-tertiary-container);
    color: var(--md-sys-color-on-tertiary-container);
  }

  .row-actions {
    display: flex;
    gap: var(--bb-space-xs);
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border: none;
    background: transparent;
    border-radius: var(--bb-shape-full);
    color: var(--md-sys-color-on-surface-variant);
    cursor: pointer;
    transition: background-color 150ms ease;
  }

  .icon-btn:hover:not(:disabled) {
    background: var(--md-sys-color-surface-container-highest);
  }

  .icon-btn.danger:hover:not(:disabled) {
    background: var(--md-sys-color-error-container);
    color: var(--md-sys-color-on-error-container);
  }

  .icon-btn:disabled {
    opacity: 0.35;
    cursor: not-allowed;
  }

  .icon-btn .material-symbols-outlined {
    font-size: 20px;
  }

  .overlay {
    position: fixed;
    inset: 0;
    z-index: 100;
    background: color-mix(in srgb, var(--md-sys-color-scrim, #000) 45%, transparent);
    backdrop-filter: blur(2px);
    display: flex;
    align-items: center;
    justify-content: center;
    padding: var(--bb-space-md);
  }

  .dialog {
    width: 100%;
    max-width: 480px;
    max-height: calc(100dvh - 48px);
    overflow-y: auto;
    padding: var(--bb-space-lg);
    border-radius: var(--bb-shape-lg);
  }

  .dialog h2 {
    margin: 0 0 var(--bb-space-md);
    color: var(--md-sys-color-on-surface);
  }

  .dialog form {
    display: flex;
    flex-direction: column;
    gap: var(--bb-space-md);
  }

  .dialog-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--bb-space-sm);
    margin-top: var(--bb-space-sm);
  }

  .text-btn {
    border: none;
    background: transparent;
    color: var(--md-sys-color-primary);
    padding: 12px 20px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
  }

  .text-btn:hover {
    background: var(--md-sys-color-surface-container-high);
  }

  .danger-btn {
    border: none;
    background: var(--md-sys-color-error);
    color: var(--md-sys-color-on-error);
    padding: 12px 24px;
    border-radius: var(--bb-shape-full);
    cursor: pointer;
  }

  .danger-btn:hover {
    opacity: 0.9;
  }
</style>
