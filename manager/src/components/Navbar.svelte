<script>
    import {link} from 'svelte-spa-router'
</script>

<nav class="navbar is-dark" role="navigation" aria-label="main navigation">
    <div class="navbar-brand">
        <a class="navbar-item" href="/" use:link>
        PIOT Manager
        </a>
    </div>

    <a class="navbar-item" href="/users" use:link>Users</a>

    <div class="navbar-end">
        <div class="navbar-item">
            <div class="buttons">
                <a href="/signup" use:link class="button is-primary"><strong>Sign up</strong></a>
                <a href="/login" use:link class="button is-light">Log in</a>
            </div>
        </div>
    </div>
</nav>
