# md2hugo

Covert your markdown files to Hugo format.

Initially designed for notes exported from [Bear](https://bear.app) but works for any markdown docs.

## Usage

Use the following command to convert your markdown documents:

```bash
md2hugo src_dir dst_dir
```

where:

- `src_dir` is a directory that contains the markdown files you want to convert which should meet some simple [requirements](#markdown-requirements) described below.
- `dst_dir` is a directory to convert your documents into. File name is preserved (i.e. `src_dir/mypost.md` will be converted into `dst_dir/mypost.md`)

## Markdown Requirements

In order to produce the best conversion results for Hugo, your source documents should follow a few requirements:

- Your markdown file must have a title as the **first line** that starts with `#` such as `# A trip to Tokyo`. You can even starts with `##` or `###` if you want to. `md2hugo` will use this line to find the title field for your output documents and any leading `#` will be removed.
- Optionally, you can include the tags for your markdown file in the **second line** in the form of multiple `#your-tag-name separated by white spaces such as `#mytag1 #mytag2`. `md2hugo` will use this line to find tags for your output documents and any leading `#` will be removed as well.

## Base tag

`md2hugo` is initially designed for markdown files written in Bear. A common situation will be that you write a lot of docs but you only want to publish only a few of them to Hugo. To do that, you can manually select the docs you want to publish, export them from Bear to a directory and use `md2hugo` to do the conversion. However, all of your tags in the source doc will also be available in your destination directory which may cause a little privacy leakage. Image you've written a doc about the comparison between different jewelries and tagged it with `#lifestyle` and #for-my-girl. Maybe you don't want to publish the second one for obvious reasons, then we need a solution here.

In Bear, documents can have [nested tags](https://bear.app/faq/Tags%20&%20Linking/Nested%20Tags/) which is a perfect a solution for this problem. You can use a tag as the parent tag for everything you want to publish (such as `#hugo`) and have any tags as the children (such as `#hugo/algorithm`, `#hugo/golang`) . On one hand you now have a good way to organize all your published documents in Bear. On the other hand, you can use the `-T` flag to specify your base tag (`hugo` in this case) and `md2hugo` will only grab tags that start with the given base tag (``#hugo/golang`), remove the leading base part (`golang`) and use them as the tags for your output documents.

## An example workflow with Bear

Here's an example of how everything works together with Bear:

1. Write a markdown doc in Bear
2. Tag it with tags start with `#hugo`
3. Export it to a directory such as src`_dir`. If you want to export all documents, click the `Hugo` tag in the left side tag bar to show all the docs. Select them all and export.
4. Create `dst_dir` and run `md2hugo -T hugo src_dir dst_dir `
5. Now you should have all the documents you want to add to Hugo. Go head using Hugo to build your website.

