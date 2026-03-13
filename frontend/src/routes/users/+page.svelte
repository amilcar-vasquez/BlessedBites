<script lang="ts">
  import { onMount } from 'svelte';
  import {
    deleteAdminUser,
    listAdminUsers,
    updateAdminUser,
    type AdminUser
  } from '$lib/api/admin';

  type EditDraft = {
    email: string;
    full_name: string;
    phone_no: string;
    role: string;
  };

  let users: AdminUser[] = [];
  let loading = true;
  let error = '';
  let savingId: number | null = null;
  let draftById: Record<number, EditDraft> = {};

  function toDraft(user: AdminUser): EditDraft {
    return {
      email: user.email || '',
      full_name: user.full_name || '',
      phone_no: user.phone_no || '',
      role: (user.role || 'customer').toLowerCase()
    };
  }

  async function loadUsers() {
    loading = true;
    error = '';
    try {
      users = await listAdminUsers();
      const nextDrafts: Record<number, EditDraft> = {};
      users.forEach((u) => {
        nextDrafts[u.id] = toDraft(u);
      });
      draftById = nextDrafts;
    } catch (e) {
      console.error(e);
      error = 'Could not load users. Ensure you are logged in as admin.';
    } finally {
      loading = false;
    }
  }

  async function saveUser(id: number) {
    const draft = draftById[id];
    if (!draft) return;
    savingId = id;
    error = '';
    try {
      const updated = await updateAdminUser(id, {
        email: draft.email.trim(),
        full_name: draft.full_name.trim(),
        phone_no: draft.phone_no.trim(),
        role: draft.role
      });

      users = users.map((user) => (user.id === id ? updated : user));
      draftById = { ...draftById, [id]: toDraft(updated) };
    } catch (e) {
      console.error(e);
      error = 'Could not save user changes.';
    } finally {
      savingId = null;
    }
  }

  function cancelEdit(id: number) {
    const current = users.find((u) => u.id === id);
    if (!current) return;
    draftById = { ...draftById, [id]: toDraft(current) };
  }

  async function removeUser(id: number) {
    if (!confirm('Delete this user account? This cannot be undone.')) return;
    savingId = id;
    error = '';
    try {
      await deleteAdminUser(id);
      users = users.filter((user) => user.id !== id);
      const { [id]: _, ...rest } = draftById;
      draftById = rest;
    } catch (e) {
      console.error(e);
      error = 'Could not delete user.';
    } finally {
      savingId = null;
    }
  }

  onMount(async () => {
    await loadUsers();
  });
</script>

<main class="shell">
  <h1>Users</h1>
  <p class="muted">Admin user management: update profile fields, role, and delete accounts.</p>

  <button class="refresh" type="button" on:click={loadUsers} disabled={loading}>Refresh</button>

  {#if loading}
    <p>Loading users...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else if users.length === 0}
    <p>No users found.</p>
  {:else}
    <div class="card table-wrap">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>Email</th>
            <th>Full Name</th>
            <th>Phone</th>
            <th>Role</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each users as user (user.id)}
            <tr>
              <td>{user.id}</td>
              <td><input bind:value={draftById[user.id].email} /></td>
              <td><input bind:value={draftById[user.id].full_name} /></td>
              <td><input bind:value={draftById[user.id].phone_no} /></td>
              <td>
                <select bind:value={draftById[user.id].role}>
                  <option value="customer">customer</option>
                  <option value="admin">admin</option>
                </select>
              </td>
              <td class="actions">
                <button class="btn" type="button" on:click={() => saveUser(user.id)} disabled={savingId === user.id}>
                  {savingId === user.id ? 'Saving...' : 'Save'}
                </button>
                <button class="btn ghost" type="button" on:click={() => cancelEdit(user.id)} disabled={savingId === user.id}>
                  Cancel
                </button>
                <button class="btn danger" type="button" on:click={() => removeUser(user.id)} disabled={savingId === user.id}>
                  Delete
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</main>

<style>
  .shell { max-width: 1000px; margin: 1.2rem auto; padding: 1rem; }
  .muted { color: #6f4744; }
  .refresh { border: none; border-radius: 999px; padding: 0.5rem 0.85rem; background: #f6dfb7; color: #492316; font-weight: 700; margin-bottom: 0.8rem; }
  .card { background: #fffaf5; border: 1px solid #e7d3c9; border-radius: 20px; padding: 0.8rem; }
  .table-wrap { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; }
  th, td { text-align: left; padding: 0.55rem 0.4rem; border-bottom: 1px solid #f0dfd6; }
  input, select { border: 1px solid #d9c5bc; border-radius: 10px; padding: 0.45rem 0.6rem; font: inherit; min-width: 140px; }
  .actions { display: flex; gap: 0.45rem; flex-wrap: wrap; }
  .btn { border: none; border-radius: 999px; padding: 0.4rem 0.7rem; background: #7f1d2d; color: #fff; font-weight: 700; }
  .btn.ghost { background: #eee0d5; color: #482826; }
  .btn.danger { background: #ffd9de; color: #5c1524; }
  .error { color: #8a1732; }
</style>
