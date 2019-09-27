<script>
    import {authenticated} from '../stores.js'
    import {link} from 'svelte-spa-router'
    import {profile} from '../stores'

</script>

<nav class="navbar is-dark" role="navigation" aria-label="main navigation">

    <div class="container">

        <div class="navbar-brand">
            <a class="navbar-item" href="/" use:link>
            PIOT Manager
            </a>

            <a role="button" class="navbar-burger burger" aria-label="menu" aria-expanded="false" data-target="navbarMenu">
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
                <span aria-hidden="true"></span>
            </a>

        </div>

        <div id="navbarMenu" class="navbar-menu">

            <div class="navbar-start">
                {#if $authenticated}
                <a class="navbar-item" href="/orgs" use:link>Organizations</a>
                <a class="navbar-item" href="/users" use:link>Users</a>
                <a class="navbar-item" href="/things" use:link>Things</a>
                {/if}
            </div>

            <div class="navbar-end">
                <div class="navbar-item">
                    <div class="buttons">
                        {#if $authenticated}
                        <a href="/signout" use:link class="button is-light">Sign out {#if $profile}{$profile.email}{/if}</a>
                        {:else}
                        <a href="/login" use:link class="button is-light">Log in</a>
                        {/if}
                    </div>
                </div>
            </div>

        </div>

    </div>

</nav>
