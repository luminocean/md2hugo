package md2hugo

import (
    "bufio"
    "fmt"
    "gopkg.in/yaml.v2"
    "os"
    "path"
    "path/filepath"
    "strings"
    "time"
)

const MarkdownExtension = ".md"

var TagBase string

func ConvertAll(srcDir string, dstDir string) error {
    if ok, err := IsDirectory(srcDir); err != nil {
        return err
    } else if !ok {
        return fmt.Errorf("%s is not a directory", srcDir)
    }
    dir, err := os.Open(srcDir)
    if err != nil {
        return err
    }
    defer dir.Close()

    entries, err := dir.Readdirnames(-1)
    if err != nil {
        return err
    }
    // handle top level regular files only
    for _, entry := range entries {
        entryPath := path.Join(srcDir, entry)
        if ok, err := IsFile(entryPath); err != nil {
            return err
        } else if !ok {
            Warnf("%s is a directory. ignore", entry)
            continue
        } else if !strings.HasSuffix(entry, MarkdownExtension) {
            continue
        }
        err = convertMarkdown(entryPath, dstDir)
        if err != nil {
            return err
        }
    }
    return nil
}

func convertMarkdown(srcFilePath string, dstDir string) error {
    src, err := os.Open(srcFilePath)
    if err != nil {
        return err
    }
    defer src.Close()

    baseName := strings.TrimSuffix(filepath.Base(srcFilePath), MarkdownExtension)
    dstFilePath := path.Join(dstDir, baseName + MarkdownExtension)
    Logf("converting %s into %s", srcFilePath, dstFilePath)

    dst, err := os.Create(dstFilePath)
    if err != nil {
        return err
    }
    defer dst.Close()

    scanner := bufio.NewScanner(src)

    // parse the first line as title
    title := ""
    ok := scanner.Scan()
    if err = scanner.Err(); !ok {
        if err != nil {
            return err
        } else {
            // EOF
            return fmt.Errorf("no title line found in the source file %s", srcFilePath)
        }
    } else {
        title = scanner.Text()
    }

    // find and parse tag lines in the format of Bear tags such as "#aa #bb #cc"
    // see https://blog.bear.app/2020/05/getting-started-with-using-and-organizing-tags-in-bear/
    // note that to prevent ambiguity, only the second line in the doc is parsed for tags
    var secondLine, tagLine string
    ok = scanner.Scan()
    if err = scanner.Err(); !ok {
        if err != nil {
            return err
        }
        // EOF
    } else {
        secondLine = scanner.Text()
        // if the second line starts with #, this line is considered as the tag line for tags
        if strings.HasPrefix(secondLine, "#") {
            tagLine = secondLine
        } else {
            // otherwise this is the first line in the doc body
            // write to dst immediately
            if strings.TrimSpace(secondLine) != "" {
                _, err = dst.WriteString(secondLine + "\n")
                if err != nil {
                    return err
                }
            }
        }
    }
    // generate hugo front matter
    yml, err := NewFrontMatter(title, tagLine).YAML()
    if err != nil {
        return err
    }
    // write frontmatter header as YAML
    fmYAMLHeader := []string{"---\n", yml, "---\n\n"}
    for _, line := range fmYAMLHeader {
        _, err = dst.WriteString(line)
        if err != nil {
            return err
        }
    }

    skipBlank := true
    // write rest of the contents
    for scanner.Scan() {
        // skip leading blank lines for better look in the result md file
        if skipBlank {
            if strings.TrimSpace(scanner.Text()) == "" {
                continue
            } else {
                skipBlank = false
            }
        }
        _, err = dst.WriteString(scanner.Text() + "\n")
        if err != nil {
            return err
        }
    }
    if err := scanner.Err(); err != nil {
        return err
    }
    return nil
}

// FrontMatter is the model for Hugo front matter
type FrontMatter struct {
    Title string   `yaml:"title"`
    Date  string   `yaml:"date"`
    Tags  []string `yaml:"tags"`
    Draft bool     `yaml:"draft"`
}

func NewFrontMatter(title, tagLine string) FrontMatter {
    // remove any leading # for the title line
    for strings.HasPrefix(title, "#") {
        title = title[1:]
    }

    tags := []string{}
    for _, segment := range strings.Split(tagLine, " ") {
        if segment == "" || !strings.HasPrefix(segment, "#") {
            continue
        }
        tag := strings.TrimPrefix(segment, "#")
        // tag base specified
        // ignore unmatched tags and trim prefix for matched tags
        if TagBase != "" {
            // ensure tailing slash
            if !strings.HasSuffix(TagBase, "/") {
                TagBase = TagBase + "/"
            }
            if !strings.HasPrefix(tag, TagBase) {
                continue
            } else {
                tag = strings.TrimPrefix(tag, TagBase)
            }
        }
        tags = append(tags, tag)
    }

    return FrontMatter{
        Title: title,
        Tags:  tags,
        Date:  time.Now().Format(time.RFC3339),
        Draft: false,
    }
}

// YAML serializes a front matter into YAML
func (m FrontMatter) YAML() (string, error) {
    bytes, err := yaml.Marshal(m)
    if err != nil {
        return "", err
    }
    return string(bytes), nil
}
