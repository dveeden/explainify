# Explainify

Explainify is a pager program for the [MySQL Client](https://dev.mysql.com/doc/en/mysql.html) that attempts to make `EXPLAIN` and other output to look nicer.

## Features

### Unicode tables

```
mysql> WITH RECURSIVE cte(n) AS (SELECT 1 n UNION ALL SELECT n+1 FROM cte WHERE n<5) SELECT n, REPEAT("abc",n) FROM cte;
+------+-----------------+
| n    | REPEAT("abc",n) |
+------+-----------------+
|    1 | abc             |
|    2 | abcabc          |
|    3 | abcabcabc       |
|    4 | abcabcabcabc    |
|    5 | abcabcabcabcabc |
+------+-----------------+
5 rows in set (0.00 sec)
```

becomes

```
mysql> WITH RECURSIVE cte(n) AS (SELECT 1 n UNION ALL SELECT n+1 FROM cte WHERE n<5) SELECT n, REPEAT("abc",n) FROM cte;
╭──────┬─────────────────╮
│ n    │ REPEAT("abc",n) │
├──────┼─────────────────┤
│    1 │ abc             │
│    2 │ abcabc          │
│    3 │ abcabcabc       │
│    4 │ abcabcabcabc    │
│    5 │ abcabcabcabcabc │
╰──────┴─────────────────╯
5 rows in set (0.00 sec)
```

This also attempts to fix the issue that wide characters like emojis mess up the layout:

```
mysql> WITH RECURSIVE cte(n) AS (SELECT 1 n UNION ALL SELECT n+1 FROM cte WHERE n<5) SELECT n, REPEAT("😀",n) FROM cte;
+------+----------------------+
| n    | REPEAT("?",n)        |
+------+----------------------+
|    1 | 😀                     |
|    2 | 😀😀                     |
|    3 | 😀😀😀                     |
|    4 | 😀😀😀😀                     |
|    5 | 😀😀😀😀😀                     |
+------+----------------------+
5 rows in set (0.00 sec)
```

```
mysql> WITH RECURSIVE cte(n) AS (SELECT 1 n UNION ALL SELECT n+1 FROM cte WHERE n<5) SELECT n, REPEAT("😀",n) FROM cte;
╭──────┬──────────────────────╮
│ n    │ REPEAT("?",n)        │
├──────┼──────────────────────┤
│    1 │ 😀                   │
│    2 │ 😀😀                 │
│    3 │ 😀😀😀               │
│    4 │ 😀😀😀😀             │
│    5 │ 😀😀😀😀😀           │
╰──────┴──────────────────────╯
5 rows in set (0.00 sec)
```

### Markdown tables

This makes it easy to copy-paste tables into markdown documentation.

```
mysql> WITH RECURSIVE cte(n) AS (SELECT 1 n UNION ALL SELECT n+1 FROM cte WHERE n<5) SELECT n, REPEAT("abc",n) FROM cte;
+------+-----------------+
| n    | REPEAT("abc",n) |
+------+-----------------+
|    1 | abc             |
|    2 | abcabc          |
|    3 | abcabcabc       |
|    4 | abcabcabcabc    |
|    5 | abcabcabcabcabc |
+------+-----------------+
5 rows in set (0.00 sec)
```

becomes

```
mysql> WITH RECURSIVE cte(n) AS (SELECT 1 n UNION ALL SELECT n+1 FROM cte WHERE n<5) SELECT n, REPEAT("abc",n) FROM cte;
| n    | REPEAT("abc",n) |
|------|-----------------|
|    1 | abc             |
|    2 | abcabc          |
|    3 | abcabcabc       |
|    4 | abcabcabcabc    |
|    5 | abcabcabcabcabc |
5 rows in set (0.00 sec)
```

And rendered this looks like:

| n    | REPEAT("abc",n) |
|------|-----------------|
|    1 | abc             |
|    2 | abcabc          |
|    3 | abcabcabc       |
|    4 | abcabcabcabc    |
|    5 | abcabcabcabcabc |

### Improved explain output for JSON and TREE

```
mysql> EXPLAIN FORMAT=json SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3\G
*************************** 1. row ***************************
EXPLAIN: {
  "query_block": {
    "union_result": {
      "using_temporary_table": false,
      "query_specifications": [
        {
          "dependent": false,
          "cacheable": true,
          "query_block": {
            "select_id": 1,
            "message": "No tables used"
          }
        },
        {
          "dependent": false,
          "cacheable": true,
          "query_block": {
            "select_id": 2,
            "message": "No tables used"
          }
        },
        {
          "dependent": false,
          "cacheable": true,
          "query_block": {
            "select_id": 3,
            "message": "No tables used"
          }
        }
      ]
    }
  }
}
1 row in set, 1 warning (0.01 sec)
```

```
mysql> EXPLAIN FORMAT=json SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3;
          EXPLAIN
--------------------------------------------------------------------------------
{
  "query_block": {
    "union_result": {
      "using_temporary_table": false,
      "query_specifications": [
        {
          "dependent": false,
          "cacheable": true,
          "query_block": {
            "select_id": 1,
            "message": "No tables used"
          }
        },
        {
          "dependent": false,
          "cacheable": true,
          "query_block": {
            "select_id": 2,
            "message": "No tables used"
          }
        },
        {
          "dependent": false,
          "cacheable": true,
          "query_block": {
            "select_id": 3,
            "message": "No tables used"
          }
        }
      ]
    }
  }
}
1 row in set, 1 warning (0.00 sec)
```

Note that `EXPLAIN FORMAT=JSON...` normally requires the use of `\G` or `--auto-vertical-output` as the default table based format looks even worse as it wasn't made for single row, single column, multi line output.

And it does syntax highlighting for JSON explain:

![alt text](image.png)

And it supports multiple themes:

![alt text](image-1.png)

And for the TREE format:

```
mysql> EXPLAIN FORMAT=tree SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3\G
*************************** 1. row ***************************
EXPLAIN: -> Append  (cost=0..0 rows=3)
    -> Stream results  (cost=0..0 rows=1)
        -> Rows fetched before execution  (cost=0..0 rows=1)
    -> Stream results  (cost=0..0 rows=1)
        -> Rows fetched before execution  (cost=0..0 rows=1)
    -> Stream results  (cost=0..0 rows=1)
        -> Rows fetched before execution  (cost=0..0 rows=1)

1 row in set (0.00 sec)
```

```
mysql> EXPLAIN FORMAT=tree SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3;
          EXPLAIN
--------------------------------------------------------------------------------
-> Append  (cost=0..0 rows=3)
    -> Stream results  (cost=0..0 rows=1)
        -> Rows fetched before execution  (cost=0..0 rows=1)
    -> Stream results  (cost=0..0 rows=1)
        -> Rows fetched before execution  (cost=0..0 rows=1)
    -> Stream results  (cost=0..0 rows=1)
        -> Rows fetched before execution  (cost=0..0 rows=1)

1 row in set (0.00 sec)
```

## Usage

When starting the MySQL Client:
```
mysql --pager=/path/to/explainify
```

When already in the MySQL Client:
```
pager /path/to/explainify
```

To disable:
```
nopager
```

Options:
```
Usage of explainify:
  -format string
    	output format: plain, markdown or unicode (default "unicode")
  -theme string
    	chroma syntax highlighting theme (default "monokailight")
```

## Related

- [Bug #77279: Use unicode to draw borders (contribution)](https://bugs.mysql.com/bug.php?id=77279)