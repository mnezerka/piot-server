<script>
    import {token, authenticated} from '../stores.js'
    import {push} from 'svelte-spa-router'
    import {gql} from '../utils.js';
    import {onMount} from 'svelte';
    import ErrorBar from '../components/ErrorBar.svelte';
    import TableButtonDelete from '../components/TableButtonDelete.svelte';

    export var params;

    let alias = '';
    let thing = null;
    let error = null;
    let fetching = false;
    let success = null;
    let orgs = [];
    let orgsAssigned = [];
    let orgAdd = '';
    let orgId = null;

    onMount(() => {
        if (!$authenticated) { push("/login"); }
        fetchThing();
    })

    async function fetchThing() {
        if (!params.id) {
            error = 'No thing specified';
            return
        }

        fetching = true;
        error = null

        try {
            let data = await gql({query: `{thing(id: "${params.id}") {id, name, alias, org {id}} orgs {id, name}}`});
            thing = data.thing;
            alias = thing.alias;
            orgId = thing.org ? thing.org.id : "";
            orgs = data.orgs;
            //orgsAssigned = user.orgs.map(o => o.id);
        } catch (e) {
            error = e;
        }
        fetching = false;
    }

    async function updateThing()
    {
        fetching = true;
        error = null;
        success = false;

        try {
            let orgIdStr = orgId === "" ? "null" : `"${orgId}"`
            let data = await gql({query: `mutation {updateThing(thing: {id: "${params.id}", alias: "${alias}", orgId: ${orgIdStr}}) {id}}`});
            success = 'Thing successfully updated'
        } catch(e) {
            error = e;
        }
        fetching = false;
    }

    async function orgSet()
    {
        console.log('set org');
    }

</script>

<style>
form { width: 24rem;}
h2 { margin-top: 2rem; }
.delete-button { text-align: right; }
</style>

<h1 class="title">Thing</h1>

<ErrorBar error={error}/>

{#if success}<div class="notification is-success has-text-centered">{success}</div>{/if}

{#if thing}

<h2 class="subtitle">Edit Thing</h2>

<form on:submit|preventDefault={updateThing}>

    <div class="field">
        <label class="label">Alias</label>
        <p class="control">
            <input bind:value={alias} class="input" placeholder="Thing alias">
        </p>
    </div>

    <div class="field">
        <label class="label">Organization</label>
        <div class="control">
            <div class="select">
                <select bind:value={orgId}>
                    <option value="">---</option>
                    {#each orgs as org}
                    <option value="{org.id}">{org.name}</option>
                    {/each}
                </select>
            </div>
        </div>
    </div>

    <div class="field">
        <p class="control">
            <button class="button is-success">Update</button>
        </p>
    </div>

</form>

{/if}
