<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { toast } from "svelte-sonner";
  import { ImportClients } from "wailsjs/go/service/ImportService";
  import { fetchClients } from "$lib/stores/clients.svelte.js";

  type ImportResult = { inserted: number; skipped: number; errors: string[] };

  let { open = $bindable(false) }: { open: boolean } = $props();

  let importing = $state(false);
  let result = $state<ImportResult | null>(null);

  $effect(() => {
    if (!open) {
      result = null;
    }
  });

  async function handleImport() {
    importing = true;
    try {
      result = await ImportClients();
      if (result.inserted > 0) {
        await fetchClients();
        toast.success(`${result.inserted} clientes importados`);
      }
    } catch (e) {
      toast.error(String(e));
    } finally {
      importing = false;
    }
  }

  function handleClose() {
    open = false;
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-w-md">
    <Dialog.Header>
      <Dialog.Title>Importar Clientes</Dialog.Title>
      <Dialog.Description>
        Seleccioná un archivo CSV con columnas: <code>nombre</code>, <code>direccion</code> (opcional), <code>zona</code>.
      </Dialog.Description>
    </Dialog.Header>

    <div class="space-y-4 py-2">
      {#if result}
        <div class="space-y-2 text-sm">
          <p><span class="font-medium">{result.inserted}</span> clientes importados</p>
          <p><span class="font-medium">{result.skipped}</span> omitidos (duplicados)</p>
          {#if result.errors.length > 0}
            <div>
              <p class="font-medium text-destructive">{result.errors.length} errores:</p>
              <ul class="mt-1 max-h-40 overflow-y-auto rounded border p-2 text-muted-foreground">
                {#each result.errors as err}
                  <li>{err}</li>
                {/each}
              </ul>
            </div>
          {/if}
        </div>
      {:else}
        <p class="text-sm text-muted-foreground">
          Seleccioná el archivo CSV para importar clientes.
        </p>
      {/if}
    </div>

    <Dialog.Footer>
      <Button variant="outline" onclick={handleClose}>Cerrar</Button>
      {#if !result}
        <Button onclick={handleImport} disabled={importing}>
          {importing ? "Importando..." : "Seleccionar archivo"}
        </Button>
      {/if}
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
