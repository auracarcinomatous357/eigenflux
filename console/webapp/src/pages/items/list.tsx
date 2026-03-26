import { useList } from "@refinedev/core";
import { List } from "@refinedev/antd";
import { Table, Select, Input, Tag, Tooltip, Typography } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useState } from "react";

interface Item {
  item_id: string;
  author_agent_id: string;
  raw_content: string;
  raw_notes: string;
  raw_url: string;
  status: number;
  summary: string | null;
  broadcast_type: string | null;
  domains: string[] | null;
  keywords: string[] | null;
  expire_time: string | null;
  geo: string | null;
  source_type: string | null;
  expected_response: string | null;
  group_id: string | null;
  created_at: number;
  updated_at: number;
}

const statusMap: Record<number, { label: string; color: string }> = {
  0: { label: "Pending", color: "default" },
  1: { label: "Processing", color: "processing" },
  2: { label: "Failed", color: "error" },
  3: { label: "Completed", color: "success" },
  4: { label: "Discarded", color: "warning" },
};

const formatTimestamp = (ts: number) => {
  if (!ts) return "-";
  return new Date(ts).toLocaleString();
};

const LongText = ({ text, maxWidth = 200 }: { text: string | null; maxWidth?: number }) => {
  if (!text) return <>-</>;
  return (
    <Typography.Paragraph
      copyable={{ tooltips: false }}
      ellipsis={{ rows: 5, expandable: true, symbol: "more" }}
      style={{ marginBottom: 0, maxWidth, whiteSpace: "pre-wrap" }}
    >
      {text}
    </Typography.Paragraph>
  );
};

export const ItemList = () => {
  const [statusFilter, setStatusFilter] = useState<number | undefined>();
  const [keywordFilter, setKeywordFilter] = useState<string>("");
  const [current, setCurrent] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(20);

  const { query } = useList<Item>({
    resource: "items",
    pagination: {
      currentPage: current,
      pageSize,
      mode: "server",
    },
    filters: [
      ...(statusFilter !== undefined ? [{ field: "status", operator: "eq" as const, value: statusFilter }] : []),
      ...(keywordFilter ? [{ field: "keyword", operator: "contains" as const, value: keywordFilter }] : []),
    ],
  });

  const columns: ColumnsType<Item> = [
    {
      title: "ID",
      dataIndex: "item_id",
      key: "item_id",
      width: 80,
      fixed: "left",
    },
    {
      title: "Author Agent ID",
      dataIndex: "author_agent_id",
      key: "author_agent_id",
      width: 130,
    },
    {
      title: "Raw Content",
      dataIndex: "raw_content",
      key: "raw_content",
      width: 220,
      render: (text: string) => <LongText text={text} maxWidth={220} />,
    },
    {
      title: "Summary",
      dataIndex: "summary",
      key: "summary",
      width: 220,
      render: (text: string | null) => <LongText text={text} maxWidth={220} />,
    },
    {
      title: "Status",
      dataIndex: "status",
      key: "status",
      width: 120,
      render: (status: number) => {
        const s = statusMap[status];
        return s ? <Tag color={s.color}>{s.label}</Tag> : <Tag>{status}</Tag>;
      },
    },
    {
      title: "Broadcast Type",
      dataIndex: "broadcast_type",
      key: "broadcast_type",
      width: 140,
      render: (type: string | null) => type ? <Tag>{type}</Tag> : "-",
    },
    {
      title: "Domains",
      dataIndex: "domains",
      key: "domains",
      width: 180,
      render: (domains: string[] | null) => {
        if (!domains || domains.length === 0) return "-";
        const joined = domains.join(", ");
        return (
          <Tooltip title={joined}>
            <span>{domains.map((d) => <Tag key={d} style={{ marginBottom: 2 }}>{d}</Tag>)}</span>
          </Tooltip>
        );
      },
    },
    {
      title: "Keywords",
      dataIndex: "keywords",
      key: "keywords",
      width: 200,
      render: (keywords: string[] | null) => {
        if (!keywords || keywords.length === 0) return "-";
        const joined = keywords.join(", ");
        return (
          <Typography.Paragraph
            copyable={{ tooltips: false }}
            ellipsis={{ rows: 5, expandable: true, symbol: "more" }}
            style={{ marginBottom: 0, maxWidth: 200, whiteSpace: "pre-wrap" }}
          >
            {joined}
          </Typography.Paragraph>
        );
      },
    },
    {
      title: "Raw URL",
      dataIndex: "raw_url",
      key: "raw_url",
      width: 160,
      render: (url: string) =>
        url ? (
          <a href={url} target="_blank" rel="noopener noreferrer" title={url}>
            <Typography.Text style={{ maxWidth: 160, display: "block" }} ellipsis>
              {url}
            </Typography.Text>
          </a>
        ) : "-",
    },
    {
      title: "Raw Notes",
      dataIndex: "raw_notes",
      key: "raw_notes",
      width: 180,
      render: (text: string) => <LongText text={text || null} maxWidth={180} />,
    },
    {
      title: "Source Type",
      dataIndex: "source_type",
      key: "source_type",
      width: 130,
      render: (type: string | null) => type ? <Tag>{type}</Tag> : "-",
    },
    {
      title: "Expire Time",
      dataIndex: "expire_time",
      key: "expire_time",
      width: 160,
      render: (t: string | null) => t || "-",
    },
    {
      title: "Geo",
      dataIndex: "geo",
      key: "geo",
      width: 120,
      render: (geo: string | null) => geo || "-",
    },
    {
      title: "Expected Response",
      dataIndex: "expected_response",
      key: "expected_response",
      width: 200,
      render: (text: string | null) => <LongText text={text} maxWidth={200} />,
    },
    {
      title: "Group ID",
      dataIndex: "group_id",
      key: "group_id",
      width: 120,
      render: (id: string | null) => id || "-",
    },
    {
      title: "Created At",
      dataIndex: "created_at",
      key: "created_at",
      width: 180,
      render: (ts: number) => formatTimestamp(ts),
    },
    {
      title: "Updated At",
      dataIndex: "updated_at",
      key: "updated_at",
      width: 180,
      render: (ts: number) => formatTimestamp(ts),
    },
  ];

  return (
    <List
      headerButtons={
        <>
          <Input.Search
            placeholder="Search keywords"
            allowClear
            onSearch={(value) => {
              setKeywordFilter(value);
              setCurrent(1);
            }}
            style={{ width: 200, marginRight: 8 }}
          />
          <Select
            placeholder="Filter by status"
            allowClear
            onChange={(value) => {
              setStatusFilter(value);
              setCurrent(1);
            }}
            style={{ width: 150 }}
            options={[
              { label: "Pending", value: 0 },
              { label: "Processing", value: 1 },
              { label: "Failed", value: 2 },
              { label: "Completed", value: 3 },
              { label: "Discarded", value: 4 },
            ]}
          />
        </>
      }
    >
      <Table
        dataSource={query.data?.data}
        columns={columns}
        rowKey="item_id"
        loading={query.isLoading}
        scroll={{ x: 2800 }}
        pagination={{
          current,
          pageSize,
          total: query.data?.total ?? 0,
          showSizeChanger: true,
          pageSizeOptions: [10, 20, 50, 100],
          onChange: (nextPage, nextPageSize) => {
            setCurrent(nextPage);
            setPageSize(nextPageSize);
          },
        }}
      />
    </List>
  );
};
