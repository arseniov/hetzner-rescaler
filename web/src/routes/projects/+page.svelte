<script lang="ts">
  import { onMount } from 'svelte';
  import {
    Alert,
    Button,
    Card,
    Input,
    Label,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from 'flowbite-svelte';
  import { m } from '$lib/paraglide/messages.js';
  import { api } from '$lib/api';
  import type { Project } from '$lib/types';

  let projects = $state<Project[]>([]);
  let error = $state<string | null>(null);
  let newName = $state('');
  let newToken = $state('');
  let busy = $state(false);

  async function refresh() {
    try {
      projects = await api.listProjects();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    }
  }

  async function create() {
    error = null;
    busy = true;
    try {
      await api.createProject({ name: newName.trim(), hcloud_token: newToken.trim() });
      newName = '';
      newToken = '';
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    } finally {
      busy = false;
    }
  }

  async function remove(p: Project) {
    if (!confirm(m.projects_delete_confirm())) return;
    try {
      await api.deleteProject(p.id);
      await refresh();
    } catch (e) {
      error = e instanceof Error ? e.message : String(e);
    }
  }

  onMount(refresh);
</script>

<div class="p-6 max-w-5xl mx-auto space-y-6">
  <h1 class="text-3xl font-semibold text-gray-900 dark:text-white">{m.projects_title()}</h1>

  {#if error}
    <Alert color="danger">{error}</Alert>
  {/if}

  <Card>
    {#if projects.length === 0}
      <p class="text-sm text-gray-600 dark:text-gray-400">{m.projects_empty()}</p>
    {:else}
      <Table hoverable>
        <TableHead>
          <TableHeadCell>Name</TableHeadCell>
          <TableHeadCell>Token</TableHeadCell>
          <TableHeadCell>{m.projects_last_error()}</TableHeadCell>
          <TableHeadCell>Updated</TableHeadCell>
          <TableHeadCell><span class="sr-only">Actions</span></TableHeadCell>
        </TableHead>
        <TableBody>
          {#each projects as p (p.id)}
            <TableBodyRow>
              <TableBodyCell>
                <a
                  href="/projects/{p.id}"
                  class="font-medium text-blue-600 dark:text-blue-400 hover:underline">{p.name}</a
                >
              </TableBodyCell>
              <TableBodyCell
                >{p.has_token ? m.projects_token_stored() : m.projects_no_token()}</TableBodyCell
              >
              <TableBodyCell class="text-xs text-gray-600 dark:text-gray-400"
                >{p.last_error ?? ''}</TableBodyCell
              >
              <TableBodyCell>{new Date(p.updated_at).toLocaleString()}</TableBodyCell>
              <TableBodyCell>
                <Button size="xs" color="danger" onclick={() => remove(p)}
                  >{m.projects_delete_submit()}</Button
                >
              </TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    {/if}
  </Card>

  <Card>
    <h3 class="text-lg font-medium mb-3 text-gray-900 dark:text-white">
      {m.projects_new_submit()}
    </h3>
    <form
      onsubmit={(e) => {
        e.preventDefault();
        create();
      }}
      class="space-y-3"
    >
      <Label>
        {m.projects_new_name_label()}
        <Input bind:value={newName} required />
      </Label>
      <Label>
        {m.projects_new_token_label()}
        <Input type="password" bind:value={newToken} required />
      </Label>
      <Button type="submit" disabled={busy}>{m.projects_new_submit()}</Button>
    </form>
  </Card>
</div>