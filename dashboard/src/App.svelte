<script>
  import { onMount } from 'svelte';

  let currentTab = $state('overview');
  let token = $state(localStorage.getItem('alisa_token') || '');
  let authenticated = $state(!!token);
  let inputToken = $state('');

  let auditLogs = $state([]);
  let loadingLogs = $state(false);

  function handleLogin() {
    if (inputToken.trim()) {
      localStorage.setItem('alisa_token', inputToken);
      token = inputToken;
      authenticated = true;
      fetchLogs();
    }
  }

  function handleLogout() {
    localStorage.removeItem('alisa_token');
    token = '';
    authenticated = false;
    auditLogs = [];
  }

  async function fetchLogs() {
    if (!authenticated) return;
    loadingLogs = true;
    try {
      const response = await fetch('/api/v1/audit?limit=25', {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (response.status === 401) {
        handleLogout();
        return;
      }
      if (response.ok) {
        auditLogs = await response.json();
      }
    } catch (err) {
      console.error('Telemetry acquisition failed:', err);
    } finally {
      loadingLogs = false;
    }
  }

  onMount(() => {
    if (authenticated) {
      fetchLogs();
    }
  });
</script>

{#if !authenticated}
  <div class="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
    <div class="sm:mx-auto w-full sm:max-w-md bg-neutral-900 border border-neutral-800 p-8 rounded-lg shadow-2xl">
      <h2 class="text-center text-xl font-mono font-bold tracking-tight text-white">ALISA SECURITY REQUISITION</h2>
      <div class="mt-6 space-y-4">
        <div>
          <label for="token" class="block text-xs font-mono uppercase text-neutral-400">System Authorization Token</label>
          <input id="token" type="password" bind:value={inputToken} class="mt-2 block w-full rounded border border-neutral-700 bg-neutral-950 px-3 py-1.5 font-mono text-sm text-white focus:border-neutral-500 focus:outline-none" />
        </div>
        <button onclick={handleLogin} class="w-full rounded bg-neutral-800 hover:bg-neutral-700 border border-neutral-700 py-2 px-4 font-mono text-sm uppercase text-white transition-colors">Verify Credentials</button>
      </div>
    </div>
  </div>
{:else}
  <div class="flex h-full">
    <aside class="w-64 border-r border-neutral-800 bg-neutral-900 flex flex-col justify-between font-mono text-sm">
      <div class="p-4">
        <div class="mb-8 px-2 py-1 bg-neutral-950 rounded border border-neutral-800">
          <span class="text-white font-bold tracking-wider">ALISA CORE</span>
          <span class="text-xs text-neutral-500 block">v1.0.0 // 512MB Target</span>
        </div>
        <nav class="space-y-1">
          <button onclick={() => currentTab = 'overview'} class="w-full text-left px-3 py-2 rounded transition-colors {currentTab === 'overview' ? 'bg-neutral-800 text-white border-l-2 border-neutral-400' : 'text-neutral-400 hover:bg-neutral-800/50'}">Overview</button>
          <button onclick={() => { currentTab = 'logging'; fetchLogs(); }} class="w-full text-left px-3 py-2 rounded transition-colors {currentTab === 'logging' ? 'bg-neutral-800 text-white border-l-2 border-neutral-400' : 'text-neutral-400 hover:bg-neutral-800/50'}">Logging Panel</button>
        </nav>
      </div>
      <div class="p-4 border-t border-neutral-800 bg-neutral-950">
        <button onclick={handleLogout} class="w-full text-center text-xs uppercase tracking-wider text-red-400 hover:text-red-300 transition-colors">Terminate Session</button>
      </div>
    </aside>

    <main class="flex-1 overflow-y-auto bg-neutral-950 p-8 font-mono">
      {#if currentTab === 'overview'}
        <header class="mb-6">
          <h1 class="text-2xl font-bold text-white tracking-tight">SYSTEM OVERVIEW</h1>
          <p class="text-xs text-neutral-500 mt-1">Real-time local cluster health diagnostics.</p>
        </header>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div class="border border-neutral-800 bg-neutral-900 p-4 rounded">
            <span class="text-xs text-neutral-400 uppercase">Operational Status</span>
            <div class="mt-2 text-xl font-bold text-emerald-400">ONLINE</div>
          </div>
          <div class="border border-neutral-800 bg-neutral-900 p-4 rounded">
            <span class="text-xs text-neutral-400 uppercase">Database Interface</span>
            <div class="mt-2 text-xl font-bold text-white">SQLite WAL</div>
          </div>
          <div class="border border-neutral-800 bg-neutral-900 p-4 rounded">
            <span class="text-xs text-neutral-400 uppercase">Memory Footprint Goal</span>
            <div class="mt-2 text-xl font-bold text-neutral-400">&lt; 512 MB</div>
          </div>
        </div>
      {:else if currentTab === 'logging'}
        <header class="mb-6 flex justify-between items-center">
          <div>
            <h1 class="text-2xl font-bold text-white tracking-tight">LOGGING STREAM</h1>
            <p class="text-xs text-neutral-500 mt-1">Direct chronological trace of active guild payloads.</p>
          </div>
          <button onclick={fetchLogs} class="rounded bg-neutral-900 border border-neutral-700 hover:bg-neutral-800 text-xs text-white px-3 py-1.5 transition-colors">Refresh Stream</button>
        </header>

        <div class="border border-neutral-800 bg-neutral-900 rounded overflow-hidden">
          <table class="w-full text-left text-xs text-neutral-300">
            <thead class="bg-neutral-950 uppercase text-neutral-400 border-b border-neutral-800">
              <tr>
                <th class="p-3">Log ID</th>
                <th class="p-3">Event Type</th>
                <th class="p-3">Target ID</th>
                <th class="p-3">Timestamp</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-neutral-800">
              {#if loadingLogs}
                <tr>
                  <td colspan="4" class="p-4 text-center text-neutral-500">Acquiring stream buffer sequence...</td>
                </tr>
              {:else if auditLogs.length === 0}
                <tr>
                  <td colspan="4" class="p-4 text-center text-neutral-500">No telemetry payloads registered in current trace index.</td>
                </tr>
              {:else}
                {#each auditLogs as logItem}
                  <tr class="hover:bg-neutral-800/30 transition-colors">
                    <td class="p-3 font-bold text-neutral-400">#{logItem.log_id}</td>
                    <td class="p-3"><span class="px-2 py-0.5 rounded bg-neutral-950 border border-neutral-800 text-neutral-200">{logItem.event_type}</span></td>
                    <td class="p-3 text-neutral-400">{logItem.user_id || 'SYSTEM'}</td>
                    <td class="p-3 text-neutral-500">{new Date(logItem.timestamp * 1000).toISOString()}</td>
                  </tr>
                {/each}
              {/if}
            </tbody>
          </table>
        </div>
      {/if}
    </main>
  </div>
{/if}