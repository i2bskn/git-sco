Git sloppy checkout
===================

Usage
-----

```
git sco <branch>
```

If specified branch exists in local branches to run the `git checkout <branch>`.  
If specified branch exists in remote branches to run the `git checkout -b <branch> origin/<branch>`.  
If specified branch does not exists in local or remote branches to run the `git checkout -b <branch>`

Installation
------------

Download from [releases](https://github.com/i2bskn/git-sco/releases) and stored in the `$PATH`.