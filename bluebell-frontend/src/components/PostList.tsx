import React, { useState, useEffect } from 'react';
import { List, Button, Card, Space, Typography, Tag } from 'antd';
import { LikeOutlined, DislikeOutlined } from '@ant-design/icons';
import { postApi } from '../services/api';
import { PostDetail } from '../types';
import { handleApiError, handleApiResponse } from '../utils/errorHandler';
import { containerStyles, cardStyles } from '../styles/common';

const { Title, Text } = Typography;

const PostList: React.FC = () => {
  const [posts, setPosts] = useState<PostDetail[]>([]);
  const [loading, setLoading] = useState(false);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);

  useEffect(() => {
    fetchPosts();
  }, [page]);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      const response = await postApi.getPostList(page);
      handleApiResponse(
        response,
        (data) => {
          setPosts(data);
          setTotal(data.length * 20); // 假设每页20条
        },
        '获取帖子列表失败'
      );
    } catch (error) {
      handleApiError(error, '获取帖子列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleVote = async (postId: number, direction: 1 | -1) => {
    try {
      const response = await postApi.votePost(postId, direction);
      if (handleApiResponse(response, () => {
        fetchPosts(); // 刷新帖子列表
      }, '投票失败')) {
        // 额外的成功处理逻辑（如果需要）
      }
    } catch (error) {
      handleApiError(error, '投票失败');
    }
  };

  const renderPostActions = (post: PostDetail) => [
    <Space key="actions">
      <Button
        type="text"
        icon={<LikeOutlined />}
        onClick={() => handleVote(post.post_id, 1)}
      >
        赞同 {post.vote_num || 0}
      </Button>
      <Button
        type="text"
        icon={<DislikeOutlined />}
        onClick={() => handleVote(post.post_id, -1)}
      >
        反对
      </Button>
    </Space>
  ];

  const renderPostMeta = (post: PostDetail) => (
    <Space size={[0, 8]} wrap>
      <Text type="secondary">作者: {post.author_name}</Text>
      <Text type="secondary">|</Text>
      <Tag color="blue">{post.community_detail.name}</Tag>
      <Text type="secondary">|</Text>
      <Text type="secondary">
        {new Date(post.create_time).toLocaleString()}
      </Text>
    </Space>
  );

  return (
    <div style={containerStyles}>
      <Card style={cardStyles}>
        <List
          loading={loading}
          itemLayout="vertical"
          pagination={{
            onChange: setPage,
            current: page,
            pageSize: 10,
            total: total,
            showSizeChanger: false
          }}
          dataSource={posts}
          renderItem={(post) => (
            <List.Item
              key={post.post_id}
              actions={renderPostActions(post)}
            >
              <List.Item.Meta
                title={<Title level={4}>{post.title}</Title>}
                description={renderPostMeta(post)}
              />
              <div style={{ marginTop: 16 }}>
                <Text>{post.content}</Text>
              </div>
            </List.Item>
          )}
        />
      </Card>
    </div>
  );
};

export default PostList; 