import React, {memo} from "react";
import ReactMarkdown from "react-markdown";
import {marked} from "marked";

function parseMarkdownIntoBlocks(markdown: string): string[] {
    const tokens = marked.lexer(markdown);
    return tokens.map(token => token.raw);
}

interface MarkdownBlockProps {
    content: string;
    onLinkClick?: (href: string) => void;
}

const MemoizedMarkdownBlock = memo(
    ({content, onLinkClick}: MarkdownBlockProps) => {
        return (
            <ReactMarkdown
                components={{
                    a: ({href, children}) => (
                        <a
                            href={href}
                            onClick={(e) => {
                                e.preventDefault();
                                onLinkClick?.(href);
                            }}
                            className="text-blue-600 underline hover:text-blue-800"
                        >
                            {children}
                        </a>
                    ),
                    img: () => null,   // block images
                    code: () => null,  // block inline code
                    pre: () => null,   // block code blocks
                    table: () => null, // block tables
                }}
            >
                {content}
            </ReactMarkdown>
        );
    },
    (prevProps, nextProps) => prevProps.content === nextProps.content
);

MemoizedMarkdownBlock.displayName = "MemoizedMarkdownBlock";

import { useMemo } from "react";

interface MemoizedMarkdownProps {
    content: string;
    id: string;
    onLinkClick?: (href: string) => void;
}

export const MemoizedMarkdown = memo(
    ({ content, id, onLinkClick }: MemoizedMarkdownProps) => {
        const blocks = useMemo(() => parseMarkdownIntoBlocks(content), [content]);

        return (
            <>
                {blocks.map((block, index) => (
                    <MemoizedMarkdownBlock
                        content={block}
                        key={`${id}-block_${index}`}
                        onLinkClick={onLinkClick}
                    />
                ))}
            </>
        );
    }
);

MemoizedMarkdown.displayName = "MemoizedMarkdown";
