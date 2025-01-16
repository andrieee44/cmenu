# CMENU {#cmenu align="center"}

[NAME](#NAME)\
[SYNOPSIS](#SYNOPSIS)\
[DESCRIPTION](#DESCRIPTION)\
[EXAMPLE](#EXAMPLE)\
[SEE ALSO](#SEE%20ALSO)\
[AUTHOR](#AUTHOR)\

------------------------------------------------------------------------

## NAME []{#NAME}

cmenu − clipboard menu wrapper

## SYNOPSIS []{#SYNOPSIS}

**cmenu** *MENU* \[*FILE*\]

## DESCRIPTION []{#DESCRIPTION}

**cmenu** is a clipboard menu wrapper, originally designed to work with
*fzf*(1). **cmenu** wraps *MENU* to choose from JSON key−value entries
of type string in *FILE* or standard input. **cmenu** pipes all keys to
*MENU* which must then output a valid key. **cmenu** then outputs the
value associated with the key.

## EXAMPLE []{#EXAMPLE}

**\$ echo '{\"key\":\"value\"}' \| cmenu fzf**

## SEE ALSO []{#SEE ALSO}

***fzf***(1)

## AUTHOR []{#AUTHOR}

andrieee44 (andrieee44@gmail.com)

------------------------------------------------------------------------
