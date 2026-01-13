--------------------------------------------  operator信息表 ---------------------------------------------
CREATE TABLE `operators` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增长',
  `name` VARCHAR(255) NOT NULL COMMENT '运算符名称',
  `category` VARCHAR(100) NOT NULL COMMENT '运算符分类',
  `scope` JSON COMMENT '适用范围数组，如["global", "project"]',
  `definition` TEXT COMMENT '运算符定义',
  `en_description` TEXT COMMENT '英文描述',
  `cn_description` TEXT COMMENT '中文描述',
  `documentation` TEXT COMMENT '官方文档链接',
  `level` VARCHAR(50) COMMENT '难度等级',
  `genius_level` VARCHAR(50) COMMENT 'Genius等级',
  `genius_quarter` VARCHAR(50) COMMENT 'Genius季度',
  
  -- 三个时间字段
  `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（精确到秒）',
  `create_date` DATE DEFAULT (CURRENT_DATE) COMMENT '创建日期（精确到日）',
  `create_month` VARCHAR(7) COMMENT '创建月份（精确到月，格式：YYYY-MM）',
  
  PRIMARY KEY (`id`),
  -- 移除 UNIQUE KEY `uk_name` (`name`)，因为名称可以重复
  
  KEY `idx_name` (`name`) COMMENT '名称查询索引',  -- 改为普通索引
  KEY `idx_category` (`category`) COMMENT '分类查询索引',
  KEY `idx_level` (`level`) COMMENT '难度等级查询索引',
  KEY `idx_create_time` (`create_time`) COMMENT '创建时间索引',
  KEY `idx_create_date` (`create_date`) COMMENT '创建日期索引',
  KEY `idx_create_month` (`create_month`) COMMENT '创建月份索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='运算符信息表';


------------------------------------------ alpha信息表 -----------------------------------------------------
CREATE TABLE `active_alpha_list` (
  `id` VARCHAR(50) NOT NULL COMMENT 'Alpha ID，如E5AknWGK、vRVdVmKz',
  `type` VARCHAR(20) NOT NULL COMMENT 'Alpha类型：SUPER/REGULAR',
  `author` VARCHAR(50) NOT NULL COMMENT '作者ID',
  
  -- Settings字段
  `instrument_type` VARCHAR(50) COMMENT '工具类型：EQUITY',
  `region` VARCHAR(50) COMMENT '地区：IND/USA',
  `universe` VARCHAR(100) COMMENT '股票池：TOP500/TOP3000',
  `delay` INT COMMENT '延迟参数',
  `decay` INT COMMENT '衰减参数',
  `neutralization` VARCHAR(50) COMMENT '中性化方式',
  `truncation` DECIMAL(5,4) COMMENT '截断值',
  `pasteurization` VARCHAR(20) COMMENT '巴氏杀菌设置',
  `unit_handling` VARCHAR(50) COMMENT '单位处理方式',
  `nan_handling` VARCHAR(50) COMMENT 'NaN处理方式',
  `selection_handling` VARCHAR(50) COMMENT '选择处理方式',
  `selection_limit` INT COMMENT '选择限制',
  `max_trade` VARCHAR(20) COMMENT '最大交易设置',
  `language` VARCHAR(50) COMMENT '语言：FASTEXPR',
  `visualization` TINYINT(1) DEFAULT 0 COMMENT '可视化设置',
  `start_date` DATE COMMENT '开始日期',
  `end_date` DATE COMMENT '结束日期',
  `component_activation` VARCHAR(50) COMMENT '组件激活',
  `test_period` VARCHAR(20) COMMENT '测试周期',
  
  -- Alpha代码内容
  `combo_code` TEXT COMMENT '组合代码（SUPER类型）',
  `combo_description` TEXT COMMENT '组合描述（SUPER类型）',
  `combo_operator_count` INT COMMENT '组合运算符数量',
  
  `selection_code` TEXT COMMENT '选择代码（SUPER类型）',
  `selection_description` TEXT COMMENT '选择描述（SUPER类型）',
  `selection_operator_count` INT COMMENT '选择运算符数量',
  
  `regular_code` TEXT COMMENT '常规代码（REGULAR类型）',
  `regular_description` TEXT COMMENT '常规描述（REGULAR类型）',
  `regular_operator_count` INT COMMENT '常规运算符数量',
  
  -- 基础信息
  `date_created` DATETIME COMMENT '创建时间',
  `date_submitted` DATETIME COMMENT '提交时间',
  `date_modified` DATETIME COMMENT '修改时间',
  `name` VARCHAR(255) COMMENT 'Alpha名称',
  `favorite` TINYINT(1) DEFAULT 0 COMMENT '是否收藏',
  `hidden` TINYINT(1) DEFAULT 0 COMMENT '是否隐藏',
  `color` VARCHAR(50) COMMENT '颜色标签',
  `category` VARCHAR(100) COMMENT '分类',
  
  -- 标签和分类
  `tags` JSON COMMENT '标签数组',
  `classifications` JSON COMMENT '分类信息数组',
  
  `grade` VARCHAR(50) COMMENT '等级',
  `stage` VARCHAR(20) NOT NULL COMMENT '阶段：OS等',
  `status` VARCHAR(20) NOT NULL COMMENT '状态：ACTIVE等',
  
  -- IS性能指标
  `is_pnl` INT COMMENT 'IS期间PNL',
  `is_book_size` INT COMMENT 'IS期间账面大小',
  `is_long_count` INT COMMENT 'IS期间多头数量',
  `is_short_count` INT COMMENT 'IS期间空头数量',
  `is_turnover` DECIMAL(10,4) COMMENT 'IS期间换手率',
  `is_returns` DECIMAL(10,4) COMMENT 'IS期间收益率',
  `is_drawdown` DECIMAL(10,4) COMMENT 'IS期间回撤',
  `is_margin` DECIMAL(10,6) COMMENT 'IS期间保证金',
  `is_sharpe` DECIMAL(10,2) COMMENT 'IS期间夏普比率',
  `is_fitness` DECIMAL(10,2) COMMENT 'IS期间适应度',
  `is_start_date` DATE COMMENT 'IS开始日期',
  `is_self_correlation` DECIMAL(10,4) COMMENT 'IS自相关',
  `is_prod_correlation` DECIMAL(10,4) COMMENT 'IS与生产相关',
  `is_checks` JSON COMMENT 'IS检查项数组',
  
  -- OS信息
  `os_start_date` DATE COMMENT 'OS开始日期',
  `os_is_sharpe_ratio` JSON COMMENT 'OS IS夏普比率',
  `os_pre_close_sharpe_ratio` JSON COMMENT 'OS前收盘夏普比率',
  `os_checks` JSON COMMENT 'OS检查项数组',
  
  -- Train性能指标（SUPER类型特有）
  `train_pnl` INT COMMENT '训练期间PNL',
  `train_book_size` INT COMMENT '训练期间账面大小',
  `train_long_count` INT COMMENT '训练期间多头数量',
  `train_short_count` INT COMMENT '训练期间空头数量',
  `train_turnover` DECIMAL(10,4) COMMENT '训练期间换手率',
  `train_returns` DECIMAL(10,4) COMMENT '训练期间收益率',
  `train_drawdown` DECIMAL(10,4) COMMENT '训练期间回撤',
  `train_margin` DECIMAL(10,6) COMMENT '训练期间保证金',
  `train_sharpe` DECIMAL(10,2) COMMENT '训练期间夏普比率',
  `train_fitness` DECIMAL(10,2) COMMENT '训练期间适应度',
  `train_start_date` DATE COMMENT '训练开始日期',
  
  -- Test性能指标（SUPER类型特有）
  `test_pnl` INT COMMENT '测试期间PNL',
  `test_book_size` INT COMMENT '测试期间账面大小',
  `test_long_count` INT COMMENT '测试期间多头数量',
  `test_short_count` INT COMMENT '测试期间空头数量',
  `test_turnover` DECIMAL(10,4) COMMENT '测试期间换手率',
  `test_returns` DECIMAL(10,4) COMMENT '测试期间收益率',
  `test_drawdown` DECIMAL(10,4) COMMENT '测试期间回撤',
  `test_margin` DECIMAL(10,6) COMMENT '测试期间保证金',
  `test_sharpe` DECIMAL(10,2) COMMENT '测试期间夏普比率',
  `test_fitness` DECIMAL(10,2) COMMENT '测试期间适应度',
  `test_start_date` DATE COMMENT '测试开始日期',
  
  -- 其他字段
  `prod` JSON COMMENT '生产数据',
  `competitions` JSON COMMENT '比赛数据',
  `themes` JSON COMMENT '主题数组',
  `pyramids` JSON COMMENT '金字塔数组',
  `pyramid_themes` JSON COMMENT '金字塔主题',
  `team` JSON COMMENT '团队信息',
  `osmosis_points` JSON COMMENT '渗透点数',
  
  -- 三个时间字段（用于记录数据更新时间）
  `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（精确到秒）',
  `create_date` DATE DEFAULT (CURRENT_DATE) COMMENT '创建日期（精确到日）',
  `create_month` VARCHAR(7) COMMENT '创建月份（精确到月，格式：YYYY-MM）',
  
  PRIMARY KEY (`id`),
  
  -- 常用查询索引
  KEY `idx_type` (`type`) COMMENT '类型查询索引',
  KEY `idx_author` (`author`) COMMENT '作者查询索引',
  KEY `idx_region` (`region`) COMMENT '地区查询索引',
  KEY `idx_universe` (`universe`) COMMENT '股票池查询索引',
  KEY `idx_status` (`status`) COMMENT '状态查询索引',
  KEY `idx_stage` (`stage`) COMMENT '阶段查询索引',
  KEY `idx_favorite` (`favorite`) COMMENT '收藏查询索引',
  KEY `idx_sharpe` (`is_sharpe`) COMMENT '夏普比率查询索引',
  KEY `idx_fitness` (`is_fitness`) COMMENT '适应度查询索引',
  KEY `idx_date_created` (`date_created`) COMMENT 'Alpha创建时间索引',
  KEY `idx_date_submitted` (`date_submitted`) COMMENT 'Alpha提交时间索引',
  KEY `idx_create_time` (`create_time`) COMMENT '记录创建时间索引',
  KEY `idx_create_date` (`create_date`) COMMENT '记录创建日期索引',
  KEY `idx_create_month` (`create_month`) COMMENT '记录创建月份索引',
  
  -- 复合索引
  KEY `idx_type_author` (`type`, `author`) COMMENT '类型和作者复合索引',
  KEY `idx_region_universe` (`region`, `universe`) COMMENT '地区和股票池复合索引',
  KEY `idx_stage_status` (`stage`, `status`) COMMENT '阶段和状态复合索引',
  KEY `idx_sharpe_fitness` (`is_sharpe`, `is_fitness`) COMMENT '夏普和适应度复合索引'
  
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='活跃Alpha列表信息表';


------------------------------------------ 顾问 weight | value_factor 数据表 -----------------------------------------------------

CREATE TABLE `weight_value_factor` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增长',
  
  -- 基本信息
  `user_id` VARCHAR(50) NOT NULL COMMENT '用户ID',
  
  -- 日期标识
  `stat_date` DATE NOT NULL COMMENT '统计日期',
  
  -- 因子相关（当天数据）
  `weight_factor` DECIMAL(5,2) COMMENT '权重因子',
  `value_factor` DECIMAL(5,2) COMMENT '价值因子',
  
  -- 变化量（相比于前有数据的一天）
  `weight_factor_change` DECIMAL(6,3) COMMENT '权重因子变化量',
  `value_factor_change` DECIMAL(6,3) COMMENT '价值因子变化量',
  
  -- 变化率（相比于前有数据的一天）
  `weight_factor_change_rate` DECIMAL(8,4) COMMENT '权重因子变化率',
  `value_factor_change_rate` DECIMAL(8,4) COMMENT '价值因子变化率',
  
  -- 其他数据
  `data_fields_used` INT DEFAULT 0 COMMENT '使用的数据字段数量',
  `submissions_count` INT DEFAULT 0 COMMENT '总提交次数',
  `super_alpha_submissions_count` INT DEFAULT 0 COMMENT '超级Alpha提交次数',
  `mean_prod_correlation` DECIMAL(5,4) COMMENT '平均产品相关性',
  `mean_self_correlation` DECIMAL(5,4) COMMENT '平均自相关性',
  `super_alpha_mean_prod_correlation` DECIMAL(5,4) COMMENT '超级Alpha平均产品相关性',
  `super_alpha_mean_self_correlation` DECIMAL(5,4) COMMENT '超级Alpha平均自相关性',
  `university` VARCHAR(255) COMMENT '大学',
  `country` VARCHAR(100) COMMENT '国家',
  `date_started` DATE COMMENT '开始日期',
  
  -- 时间字段
  `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_date` (`user_id`, `stat_date`) COMMENT '用户和日期唯一索引',
  
  -- 查询索引
  KEY `idx_user_id` (`user_id`) COMMENT '用户ID索引',
  KEY `idx_stat_date` (`stat_date`) COMMENT '统计日期索引',
  KEY `idx_weight_factor` (`weight_factor`) COMMENT '权重因子索引',
  KEY `idx_value_factor` (`value_factor`) COMMENT '价值因子索引',
  KEY `idx_weight_factor_change` (`weight_factor_change`) COMMENT '权重因子变化量索引',
  KEY `idx_value_factor_change` (`value_factor_change`) COMMENT '价值因子变化量索引',
  KEY `idx_create_time` (`create_time`) COMMENT '创建时间索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='研究顾问wf|vf每日统计表';


------------------------------------------------- 金字塔优先推塔表 -------------------------------------------

CREATE TABLE `pyramid_alphas` (
  `id` INT NOT NULL AUTO_INCREMENT COMMENT '主键ID，自增长',
  `user_id` VARCHAR(100) NOT NULL COMMENT '用户标识，关联用户表或记录来源',
  `category_id` VARCHAR(50) COMMENT '金字塔分类ID',
  `category_name` VARCHAR(255) COMMENT '金字塔分类名称',
  `region` VARCHAR(100) COMMENT '区域',
  `delay` INT DEFAULT 0 COMMENT '延迟天数',
  `alpha_count` INT DEFAULT 0 COMMENT 'Alpha数量',
  
  -- 新增：季度标签字段
  `quarter_tag` VARCHAR(10) COMMENT '季度标签，格式如2025-Q3',
  
  -- 三个时间字段（与参考表保持一致）
  `create_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间（精确到秒）',
  `create_date` DATE DEFAULT (CURRENT_DATE) COMMENT '创建日期（精确到日）',
  `create_month` VARCHAR(7) GENERATED ALWAYS AS (DATE_FORMAT(create_time, '%Y-%m')) STORED COMMENT '创建月份（格式：YYYY-MM）',
  
  -- 查询时间范围字段（可根据API参数动态记录）
  `stat_start_date` DATE COMMENT '统计开始日期',
  `stat_end_date` DATE COMMENT '统计结束日期',
  
  PRIMARY KEY (`id`),
  
  -- 索引设计
  KEY `idx_user_id` (`user_id`) COMMENT '用户查询索引',
  KEY `idx_category_id` (`category_id`) COMMENT '分类ID索引',
  KEY `idx_region` (`region`) COMMENT '区域索引',
  KEY `idx_quarter_tag` (`quarter_tag`) COMMENT '季度标签索引',
  KEY `idx_user_quarter` (`user_id`, `quarter_tag`) COMMENT '用户+季度复合索引',
  KEY `idx_create_time` (`create_time`) COMMENT '创建时间索引',
  KEY `idx_create_date` (`create_date`) COMMENT '创建日期索引',
  KEY `idx_create_month` (`create_month`) COMMENT '创建月份索引',
  KEY `idx_stat_date_range` (`stat_start_date`, `stat_end_date`) COMMENT '统计日期范围索引'
  
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户金字塔Alpha活动统计表';


------------------------------------------------------- alpha的Combin表现表 ----------------------------------------------

CREATE TABLE combined_alpha_performance (
    id INT PRIMARY KEY AUTO_INCREMENT,
    value_factor DECIMAL(12, 6) NULL COMMENT '因子价值',
    combined_alpha_performance DECIMAL(10, 4) NULL COMMENT '综合Alpha表现',
    combined_selected_alpha_performance DECIMAL(10, 4) NULL COMMENT '综合精选Alpha表现',
    combined_power_pool_alpha_performance DECIMAL(10, 4) NULL COMMENT '综合动力池Alpha表现',
    genius_level VARCHAR(20) NULL COMMENT 'Genius等级',
    calculation_date DATE NOT NULL COMMENT '计算日期',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Alpha Combin表现指标表';