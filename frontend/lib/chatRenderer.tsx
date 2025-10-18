import {unified} from "unified";
import remarkParse from "remark-parse";
import remarkRehype from "remark-rehype";
import rehypeReact from "rehype-react";
import {createElement, Fragment, ReactNode} from "react";

export function parseMessage(text: string): ReactNode {
    const processor = unified()
        .use(remarkParse)
        .use(remarkRehype)
        .use(rehypeReact, {
            createElement: createElement,
            Fragment,
            jsx: createElement,
            jsxs: createElement,
            components: {
                a: ({href, children}: any) => (
                    <a
                        href={href}
                        // onClick={(e) => handleLinkClick(href, e)}
                        className="text-blue-600 underline hover:text-blue-800"
                    >
                        {children}
                    </a>
                ),
                // Block unwanted elements
                img: () => null,
                table: () => null,
                code: () => null,
                pre: () => null,
            },
        });

    return processor.processSync(text).result as ReactNode;
}