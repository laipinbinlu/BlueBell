import React, { useEffect, useState } from 'react';
import { List, Card, Button, message, Space, Skeleton } from 'antd';
import { LikeOutlined, DislikeOutlined } from '@ant-design/icons';
import { PostDetail } from '../types';
import { postApi } from '../services/api';

const PostList: React.FC = () => {
  const [posts, setPosts] = useState<PostDetail[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      const response = await postApi.getPostList(page);
      if (response.code === 0) {
        setPosts(response.data);
        // 假设总数是当前页数据的20倍，实际项目中应该从后端获取
        setTotal(response.data.length * 20);
      } else {
        message.error(response.message || '获取帖子列表失败');
      }
    } catch (error) {
      console.error('获取帖子列表失败:', error);
      message.error('获取帖子列表失败，请稍后重试');
    } finally {
      setLoading(false);
    }
  };

  const handleVote = async (postId: number, direction: 1 | -1) => {
    try {
      const response = await postApi.votePost(postId, direction);
      if (response.code === 0) {
        message.success(direction === 1 ? '点赞成功' : '取消点赞成功');
        fetchPosts();
      } else {
        message.error(response.message || '操作失败');
      }
    } catch (error) {
      console.error('投票失败:', error);
      message.error('操作失败，请稍后重试');
    }
  };

  useEffect(() => {
    fetchPosts();
  }, [page]);

  const renderItem = (post: PostDetail) => (
    <List.Item
      key={post.post_id}
      actions={[
        <Space>
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
      ]}
    >
      <Skeleton loading={loading} active avatar>
        <List.Item.Meta
          title={post.title}
          description={`作者: ${post.author_name} | 发布于: ${new Date(post.create_time).toLocaleString()} | 社区: ${post.community_detail.name}`}
        />
        <div style={{ marginTop: 16 }}>{post.content}</div>
      </Skeleton>
    </List.Item>
  );

  return (
    <List
      itemLayout="vertical"
      size="large"
      pagination={{
        onChange: setPage,
        current: page,
        pageSize: 10,
        total: total,
        showSizeChanger: false
      }}
      dataSource={posts}
      renderItem={renderItem}
      style={{ background: '#fff', padding: '24px' }}
    />
  );
};

export default PostList; 